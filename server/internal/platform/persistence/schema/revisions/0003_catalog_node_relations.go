package revisions

import (
	"context"
	"strings"

	"github.com/perfect-panel/server/internal/platform/persistence/catalog"
	"github.com/perfect-panel/server/internal/platform/persistence/node"
	"github.com/perfect-panel/server/internal/platform/persistence/schema"
	"github.com/perfect-panel/server/internal/platform/persistence/subscribe"
	"github.com/perfect-panel/server/internal/platform/support/tool"
	"gorm.io/gorm"
)

type catalogNodeRelationsRevision struct{}

func (catalogNodeRelationsRevision) Name() string {
	return schema.RevisionName(3, "catalog_node_relations")
}

func (catalogNodeRelationsRevision) Up(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&catalog.NodeGroup{},
		&catalog.NodeGroupNode{},
		&catalog.PlanNodeGroupRule{},
		&node.SubscriptionNodeAssignment{},
	); err != nil {
		return err
	}

	repo := catalog.NewRepository(db)
	assignments := node.NewAssignmentRepository(db)
	ctx := context.Background()

	return db.Transaction(func(tx *gorm.DB) error {
		if tx.Migrator().HasTable(&node.Node{}) {
			var nodes []node.Node
			if err := tx.Find(&nodes).Error; err != nil {
				return err
			}
			for _, item := range nodes {
				if err := repo.SyncNodeTagMemberships(ctx, item.Id, splitCSV(item.Tags), tx); err != nil {
					return err
				}
			}
		}

		if tx.Migrator().HasTable(&subscribe.Subscribe{}) {
			var subscribes []subscribe.Subscribe
			if err := tx.Find(&subscribes).Error; err != nil {
				return err
			}
			for _, item := range subscribes {
				if err := repo.SyncSubscribeSelectors(ctx, item.Id, tool.StringToInt64Slice(item.Nodes), splitCSV(item.NodeTags), tx); err != nil {
					return err
				}
			}
			for _, item := range subscribes {
				if err := assignments.RefreshAssignmentsForSubscribe(ctx, item.Id, tx); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func splitCSV(raw string) []string {
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		result = append(result, part)
	}
	return result
}
