package revisions

import (
	"github.com/perfect-panel/server/internal/platform/persistence/node"
	"github.com/perfect-panel/server/internal/platform/persistence/schema"
	"github.com/perfect-panel/server/internal/platform/persistence/system"
	"gorm.io/gorm"
)

type asyncTrustAndUsageRevision struct{}

func (asyncTrustAndUsageRevision) Name() string {
	return schema.RevisionName(5, "async_trust_and_usage")
}

func (asyncTrustAndUsageRevision) Up(db *gorm.DB) error {
	return db.AutoMigrate(
		&system.ExternalTrustEvent{},
		&node.NodeUsageReport{},
	)
}
