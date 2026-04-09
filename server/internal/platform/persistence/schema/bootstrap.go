package schema

import (
	"errors"

	"github.com/perfect-panel/server/internal/platform/persistence/auth"
	"github.com/perfect-panel/server/internal/platform/persistence/system"
	"github.com/perfect-panel/server/internal/platform/persistence/user"
	"gorm.io/gorm"
)

var ErrResetRequiresManualCleanup = errors.New("db reset requires manual cleanup once additional revisions exist")

func Bootstrap(db *gorm.DB, source string) error {
	if err := ValidateRevisionSource(source); err != nil {
		return err
	}
	if err := ensureRegistry(db); err != nil {
		return err
	}
	revision, ok := findRevision(BaselineRevisionName)
	if !ok {
		return ErrRevisionNotFound
	}
	return applyRevision(db, revision, source)
}

func Reset(db *gorm.DB, source string) error {
	if err := ValidateRevisionSource(source); err != nil {
		return err
	}
	if db == nil {
		return errors.New("db is nil")
	}
	if len(RegisteredRevisions()) > 1 {
		return ErrResetRequiresManualCleanup
	}

	if err := db.Migrator().DropTable(
		&user.AuthMethods{},
		&user.User{},
		&auth.Auth{},
		&system.System{},
		&Registry{},
	); err != nil {
		return err
	}
	return Bootstrap(db, source)
}
