package node

import (
	"context"
	"sort"

	"github.com/perfect-panel/server/internal/platform/persistence/catalog"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	catalogRevisionName  = "0003_catalog_node_relations"
	schemaRegistryTable  = "schema_registry"
	revisionStateApplied = "applied"
)

type SubscriptionNodeAssignment struct {
	Id              int64 `gorm:"primaryKey"`
	UserSubscribeId int64 `gorm:"not null;uniqueIndex:idx_subscription_node_assignments_unique;index:idx_subscription_node_assignments_user_subscribe"`
	SubscribeId     int64 `gorm:"not null;index:idx_subscription_node_assignments_subscribe"`
	NodeId          int64 `gorm:"not null;uniqueIndex:idx_subscription_node_assignments_unique;index:idx_subscription_node_assignments_node"`
}

func (SubscriptionNodeAssignment) TableName() string {
	return "subscription_node_assignments"
}

type AssignmentRepository struct {
	db      *gorm.DB
	catalog *catalog.Repository
}

func NewAssignmentRepository(db *gorm.DB) *AssignmentRepository {
	return &AssignmentRepository{
		db:      db,
		catalog: catalog.NewRepository(db),
	}
}

func (r *AssignmentRepository) Available(conn ...*gorm.DB) bool {
	db := r.conn(nil, conn...)
	if db == nil || !r.revisionApplied(db) {
		return false
	}
	return r.Installed(db)
}

func (r *AssignmentRepository) Installed(conn ...*gorm.DB) bool {
	db := r.conn(nil, conn...)
	if db == nil {
		return false
	}
	return r.catalog.Installed(db) && db.Migrator().HasTable(&SubscriptionNodeAssignment{})
}

func (r *AssignmentRepository) ListAssignedUserSubscriptionIDs(ctx context.Context, serverId int64, protocol string, tx ...*gorm.DB) ([]int64, error) {
	db := r.conn(ctx, tx...)
	if db == nil || !r.Installed(db) {
		return nil, nil
	}

	var userSubscribeIds []int64
	query := db.Model(&SubscriptionNodeAssignment{}).
		Distinct("subscription_node_assignments.user_subscribe_id").
		Joins("JOIN nodes ON nodes.id = subscription_node_assignments.node_id").
		Where("nodes.server_id = ?", serverId)
	if protocol != "" {
		query = query.Where("nodes.protocol = ?", protocol)
	}
	err := query.Where("nodes.enabled = ?", true).
		Pluck("subscription_node_assignments.user_subscribe_id", &userSubscribeIds).Error
	if err != nil {
		return nil, err
	}
	sort.Slice(userSubscribeIds, func(i, j int) bool { return userSubscribeIds[i] < userSubscribeIds[j] })
	return userSubscribeIds, nil
}

func (r *AssignmentRepository) RefreshAssignmentsForSubscribe(ctx context.Context, subscribeId int64, tx ...*gorm.DB) error {
	db := r.conn(ctx, tx...)
	if db == nil || !r.Installed(db) {
		return nil
	}

	nodeIds, err := r.catalog.ResolveNodeIDsForSubscribe(ctx, subscribeId, db)
	if err != nil {
		return err
	}
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("subscribe_id = ?", subscribeId).Delete(&SubscriptionNodeAssignment{}).Error; err != nil {
			return err
		}
		if len(nodeIds) == 0 {
			return nil
		}
		var userSubscriptions []struct {
			Id     int64
			Status uint8
		}
		if err := tx.Table("user_subscribe").
			Where("subscribe_id = ? AND status IN ?", subscribeId, []int64{0, 1}).
			Find(&userSubscriptions).Error; err != nil {
			return err
		}
		for _, userSubscription := range userSubscriptions {
			for _, nodeId := range nodeIds {
				row := SubscriptionNodeAssignment{
					UserSubscribeId: userSubscription.Id,
					SubscribeId:     subscribeId,
					NodeId:          nodeId,
				}
				if err := tx.Clauses(clause.OnConflict{
					Columns:   []clause.Column{{Name: "user_subscribe_id"}, {Name: "node_id"}},
					DoNothing: true,
				}).Create(&row).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (r *AssignmentRepository) SyncUserSubscription(ctx context.Context, userSubscribeId, subscribeId int64, status uint8, tx ...*gorm.DB) error {
	db := r.conn(ctx, tx...)
	if db == nil || !r.Installed(db) {
		return nil
	}

	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_subscribe_id = ?", userSubscribeId).Delete(&SubscriptionNodeAssignment{}).Error; err != nil {
			return err
		}
		if status != 0 && status != 1 {
			return nil
		}
		nodeIds, err := r.catalog.ResolveNodeIDsForSubscribe(ctx, subscribeId, tx)
		if err != nil {
			return err
		}
		for _, nodeId := range nodeIds {
			row := SubscriptionNodeAssignment{
				UserSubscribeId: userSubscribeId,
				SubscribeId:     subscribeId,
				NodeId:          nodeId,
			}
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "user_subscribe_id"}, {Name: "node_id"}},
				DoNothing: true,
			}).Create(&row).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *AssignmentRepository) DeleteUserSubscription(ctx context.Context, userSubscribeId int64, tx ...*gorm.DB) error {
	db := r.conn(ctx, tx...)
	if db == nil || !r.Installed(db) {
		return nil
	}
	return db.Where("user_subscribe_id = ?", userSubscribeId).Delete(&SubscriptionNodeAssignment{}).Error
}

func (r *AssignmentRepository) RefreshAssignmentsForNode(ctx context.Context, nodeId int64, tx ...*gorm.DB) error {
	db := r.conn(ctx, tx...)
	if db == nil || !r.Installed(db) {
		return nil
	}

	subscribeIds, err := r.catalog.ListSubscribeIDsForNode(ctx, nodeId, db)
	if err != nil {
		return err
	}
	for _, subscribeId := range subscribeIds {
		if err := r.RefreshAssignmentsForSubscribe(ctx, subscribeId, db); err != nil {
			return err
		}
	}
	return nil
}

func (r *AssignmentRepository) conn(ctx context.Context, tx ...*gorm.DB) *gorm.DB {
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

func (r *AssignmentRepository) revisionApplied(db *gorm.DB) bool {
	if db == nil || !db.Migrator().HasTable(schemaRegistryTable) {
		return false
	}
	var count int64
	if err := db.Table(schemaRegistryTable).
		Where("id = ? AND state = ?", catalogRevisionName, revisionStateApplied).
		Count(&count).Error; err != nil {
		return false
	}
	return count > 0
}
