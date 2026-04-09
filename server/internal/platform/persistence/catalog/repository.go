package catalog

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	schemaRegistryTable  = "schema_registry"
	revisionStateApplied = "applied"
)

type Repository struct {
	db *gorm.DB
}

type SelectorSnapshot struct {
	NodeIds []int64
	Tags    []string
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Available(conn ...*gorm.DB) bool {
	db := r.conn(nil, conn...)
	if db == nil || !r.revisionApplied(db) {
		return false
	}
	return r.Installed(db)
}

func (r *Repository) Installed(conn ...*gorm.DB) bool {
	db := r.conn(nil, conn...)
	if db == nil {
		return false
	}
	return db.Migrator().HasTable(&NodeGroup{}) &&
		db.Migrator().HasTable(&NodeGroupNode{}) &&
		db.Migrator().HasTable(&PlanNodeGroupRule{})
}

func (r *Repository) SyncNodeTagMemberships(ctx context.Context, nodeId int64, tags []string, tx ...*gorm.DB) error {
	db := r.conn(ctx, tx...)
	if db == nil || !r.Installed(db) {
		return nil
	}

	tags = normalizeTags(tags)
	return db.Transaction(func(tx *gorm.DB) error {
		var tagGroupIds []int64
		if err := tx.Model(&NodeGroup{}).
			Where("type = ?", nodeGroupTypeTag).
			Pluck("id", &tagGroupIds).Error; err != nil {
			return err
		}
		if len(tagGroupIds) > 0 {
			if err := tx.Where("node_id = ? AND node_group_id IN ?", nodeId, tagGroupIds).
				Delete(&NodeGroupNode{}).Error; err != nil {
				return err
			}
		}
		for _, tag := range tags {
			group, err := r.ensureGroup(tx, tagGroupCode(tag), tag, nodeGroupTypeTag)
			if err != nil {
				return err
			}
			row := NodeGroupNode{
				NodeGroupId: group.Id,
				NodeId:      nodeId,
			}
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "node_group_id"}, {Name: "node_id"}},
				DoNothing: true,
			}).Create(&row).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *Repository) SyncSubscribeSelectors(ctx context.Context, subscribeId int64, nodeIds []int64, tags []string, tx ...*gorm.DB) error {
	db := r.conn(ctx, tx...)
	if db == nil || !r.Installed(db) {
		return nil
	}

	nodeIds = normalizeNodeIds(nodeIds)
	tags = normalizeTags(tags)

	return db.Transaction(func(tx *gorm.DB) error {
		var existingRules []PlanNodeGroupRule
		if err := tx.Where("subscribe_id = ?", subscribeId).Find(&existingRules).Error; err != nil {
			return err
		}

		explicitCode := explicitGroupCode(subscribeId)
		explicitGroup, err := r.ensureGroup(tx, explicitCode, fmt.Sprintf("subscribe:%d:nodes", subscribeId), nodeGroupTypeExplicit)
		if err != nil {
			return err
		}
		if err := tx.Where("node_group_id = ?", explicitGroup.Id).Delete(&NodeGroupNode{}).Error; err != nil {
			return err
		}
		for _, nodeId := range nodeIds {
			row := NodeGroupNode{
				NodeGroupId: explicitGroup.Id,
				NodeId:      nodeId,
			}
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "node_group_id"}, {Name: "node_id"}},
				DoNothing: true,
			}).Create(&row).Error; err != nil {
				return err
			}
		}

		desired := make(map[int64]PlanNodeGroupRule)
		if len(nodeIds) > 0 {
			desired[explicitGroup.Id] = PlanNodeGroupRule{
				SubscribeId: subscribeId,
				NodeGroupId: explicitGroup.Id,
				RuleType:    nodeGroupTypeExplicit,
				RuleValue:   explicitCode,
			}
		}
		for _, tag := range tags {
			group, err := r.ensureGroup(tx, tagGroupCode(tag), tag, nodeGroupTypeTag)
			if err != nil {
				return err
			}
			desired[group.Id] = PlanNodeGroupRule{
				SubscribeId: subscribeId,
				NodeGroupId: group.Id,
				RuleType:    nodeGroupTypeTag,
				RuleValue:   tag,
			}
		}

		for _, existing := range existingRules {
			if _, keep := desired[existing.NodeGroupId]; keep {
				delete(desired, existing.NodeGroupId)
				continue
			}
			if err := tx.Delete(&existing).Error; err != nil {
				return err
			}
		}
		for _, rule := range desired {
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "subscribe_id"}, {Name: "node_group_id"}},
				UpdateAll: true,
			}).Create(&rule).Error; err != nil {
				return err
			}
		}
		if len(nodeIds) == 0 {
			if err := tx.Where("subscribe_id = ? AND node_group_id = ?", subscribeId, explicitGroup.Id).
				Delete(&PlanNodeGroupRule{}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *Repository) DeleteSubscribeSelectors(ctx context.Context, subscribeId int64, tx ...*gorm.DB) error {
	db := r.conn(ctx, tx...)
	if db == nil || !r.Installed(db) {
		return nil
	}

	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("subscribe_id = ?", subscribeId).Delete(&PlanNodeGroupRule{}).Error; err != nil {
			return err
		}
		var group NodeGroup
		err := tx.Where("code = ?", explicitGroupCode(subscribeId)).First(&group).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil
			}
			return err
		}
		if err := tx.Where("node_group_id = ?", group.Id).Delete(&NodeGroupNode{}).Error; err != nil {
			return err
		}
		return tx.Delete(&group).Error
	})
}

func (r *Repository) DeleteNodeMemberships(ctx context.Context, nodeId int64, tx ...*gorm.DB) error {
	db := r.conn(ctx, tx...)
	if db == nil || !r.Installed(db) {
		return nil
	}
	return db.Where("node_id = ?", nodeId).Delete(&NodeGroupNode{}).Error
}

func (r *Repository) ResolveSubscribeIDsBySelectors(ctx context.Context, nodeIds []int64, tags []string, tx ...*gorm.DB) ([]int64, error) {
	db := r.conn(ctx, tx...)
	if db == nil || !r.Installed(db) {
		return nil, nil
	}

	nodeIds = normalizeNodeIds(nodeIds)
	tags = normalizeTags(tags)

	merge := func(target map[int64]struct{}, values []int64) map[int64]struct{} {
		if target == nil {
			target = make(map[int64]struct{}, len(values))
			for _, value := range values {
				target[value] = struct{}{}
			}
			return target
		}
		next := make(map[int64]struct{})
		for _, value := range values {
			if _, ok := target[value]; ok {
				next[value] = struct{}{}
			}
		}
		return next
	}

	var result map[int64]struct{}
	if len(nodeIds) > 0 {
		var byNodes []int64
		err := db.Model(&PlanNodeGroupRule{}).
			Distinct("plan_node_group_rules.subscribe_id").
			Joins("JOIN node_group_nodes ON node_group_nodes.node_group_id = plan_node_group_rules.node_group_id").
			Where("node_group_nodes.node_id IN ?", nodeIds).
			Pluck("plan_node_group_rules.subscribe_id", &byNodes).Error
		if err != nil {
			return nil, err
		}
		result = merge(result, byNodes)
	}
	if len(tags) > 0 {
		var byTags []int64
		err := db.Model(&PlanNodeGroupRule{}).
			Distinct("subscribe_id").
			Where("rule_type = ? AND rule_value IN ?", nodeGroupTypeTag, tags).
			Pluck("subscribe_id", &byTags).Error
		if err != nil {
			return nil, err
		}
		result = merge(result, byTags)
	}
	if result == nil {
		return nil, nil
	}

	ids := make([]int64, 0, len(result))
	for id := range result {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	return ids, nil
}

func (r *Repository) ResolveNodeIDsForSubscribe(ctx context.Context, subscribeId int64, tx ...*gorm.DB) ([]int64, error) {
	db := r.conn(ctx, tx...)
	if db == nil || !r.Installed(db) {
		return nil, nil
	}

	var nodeIds []int64
	err := db.Model(&PlanNodeGroupRule{}).
		Distinct("node_group_nodes.node_id").
		Joins("JOIN node_group_nodes ON node_group_nodes.node_group_id = plan_node_group_rules.node_group_id").
		Where("plan_node_group_rules.subscribe_id = ?", subscribeId).
		Pluck("node_group_nodes.node_id", &nodeIds).Error
	if err != nil {
		return nil, err
	}
	sort.Slice(nodeIds, func(i, j int) bool { return nodeIds[i] < nodeIds[j] })
	return nodeIds, nil
}

func (r *Repository) ListSubscribeIDsForNode(ctx context.Context, nodeId int64, tx ...*gorm.DB) ([]int64, error) {
	db := r.conn(ctx, tx...)
	if db == nil || !r.Installed(db) {
		return nil, nil
	}

	var subscribeIds []int64
	err := db.Model(&PlanNodeGroupRule{}).
		Distinct("plan_node_group_rules.subscribe_id").
		Joins("JOIN node_group_nodes ON node_group_nodes.node_group_id = plan_node_group_rules.node_group_id").
		Where("node_group_nodes.node_id = ?", nodeId).
		Pluck("plan_node_group_rules.subscribe_id", &subscribeIds).Error
	if err != nil {
		return nil, err
	}
	sort.Slice(subscribeIds, func(i, j int) bool { return subscribeIds[i] < subscribeIds[j] })
	return subscribeIds, nil
}

func (r *Repository) LoadSelectorSnapshot(ctx context.Context, subscribeId int64, tx ...*gorm.DB) (*SelectorSnapshot, error) {
	db := r.conn(ctx, tx...)
	if db == nil || !r.Installed(db) {
		return &SelectorSnapshot{}, nil
	}

	var rules []PlanNodeGroupRule
	if err := db.Where("subscribe_id = ?", subscribeId).Find(&rules).Error; err != nil {
		return nil, err
	}
	snapshot := &SelectorSnapshot{}
	for _, rule := range rules {
		switch rule.RuleType {
		case nodeGroupTypeTag:
			snapshot.Tags = append(snapshot.Tags, rule.RuleValue)
		case nodeGroupTypeExplicit:
			var nodeIds []int64
			if err := db.Model(&NodeGroupNode{}).
				Where("node_group_id = ?", rule.NodeGroupId).
				Pluck("node_id", &nodeIds).Error; err != nil {
				return nil, err
			}
			snapshot.NodeIds = append(snapshot.NodeIds, nodeIds...)
		}
	}
	snapshot.NodeIds = normalizeNodeIds(snapshot.NodeIds)
	snapshot.Tags = normalizeTags(snapshot.Tags)
	return snapshot, nil
}

func (r *Repository) ensureGroup(db *gorm.DB, code, name, groupType string) (*NodeGroup, error) {
	row := NodeGroup{
		Code: code,
		Name: name,
		Type: groupType,
	}
	if err := db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "code"}},
		DoUpdates: clause.Assignments(map[string]any{
			"name": name,
			"type": groupType,
		}),
	}).Create(&row).Error; err != nil {
		return nil, err
	}
	if err := db.Where("code = ?", code).First(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repository) conn(ctx context.Context, tx ...*gorm.DB) *gorm.DB {
	if len(tx) > 0 && tx[0] != nil {
		if ctx != nil {
			return tx[0].WithContext(ctx)
		}
		return tx[0]
	}
	if r.db == nil {
		return nil
	}
	if ctx != nil {
		return r.db.WithContext(ctx)
	}
	return r.db
}

func (r *Repository) revisionApplied(db *gorm.DB) bool {
	if db == nil || !db.Migrator().HasTable(schemaRegistryTable) {
		return false
	}
	var count int64
	if err := db.Table(schemaRegistryTable).
		Where("id = ? AND state = ?", RevisionName, revisionStateApplied).
		Count(&count).Error; err != nil {
		return false
	}
	return count > 0
}

func normalizeTags(tags []string) []string {
	seen := make(map[string]struct{}, len(tags))
	out := make([]string, 0, len(tags))
	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		if tag == "" {
			continue
		}
		if _, ok := seen[tag]; ok {
			continue
		}
		seen[tag] = struct{}{}
		out = append(out, tag)
	}
	sort.Strings(out)
	return out
}

func normalizeNodeIds(nodeIds []int64) []int64 {
	seen := make(map[int64]struct{}, len(nodeIds))
	out := make([]int64, 0, len(nodeIds))
	for _, nodeId := range nodeIds {
		if nodeId <= 0 {
			continue
		}
		if _, ok := seen[nodeId]; ok {
			continue
		}
		seen[nodeId] = struct{}{}
		out = append(out, nodeId)
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

func tagGroupCode(tag string) string {
	return tagGroupCodePrefix + tag
}

func explicitGroupCode(subscribeId int64) string {
	return fmt.Sprintf("subscribe:%d:nodes", subscribeId)
}
