package revisions

import (
	"sync"

	"github.com/perfect-panel/server/internal/platform/persistence/auth"
	"github.com/perfect-panel/server/internal/platform/persistence/schema"
	"github.com/perfect-panel/server/internal/platform/persistence/system"
	"github.com/perfect-panel/server/internal/platform/persistence/user"
	"gorm.io/gorm"
)

var registerEmbeddedOnce sync.Once

func RegisterEmbedded() {
	registerEmbeddedOnce.Do(func() {
		schema.RegisterRevision(baselineRevision{})
		schema.RegisterRevision(identitySystemRevision{})
	})
}

type baselineRevision struct{}

func (baselineRevision) Name() string {
	return schema.BaselineRevisionName
}

func (baselineRevision) Up(db *gorm.DB) error {
	return db.AutoMigrate(
		&auth.Auth{},
		&system.System{},
		&user.User{},
		&user.AuthMethods{},
	)
}
