package revisions

import (
	"github.com/perfect-panel/server/internal/platform/persistence/client"
	"github.com/perfect-panel/server/internal/platform/persistence/schema"
	"gorm.io/gorm"
)

type clientApplicationsRevision struct{}

func (clientApplicationsRevision) Name() string {
	return schema.RevisionName(7, "client_applications")
}

func (clientApplicationsRevision) Up(db *gorm.DB) error {
	return db.AutoMigrate(&client.SubscribeApplication{})
}
