package system

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	modelsystem "github.com/perfect-panel/server/internal/platform/persistence/system"
	"github.com/perfect-panel/server/internal/platform/support/xerr"
	"github.com/stretchr/testify/require"
	sqliteDriver "gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type fakeSystemModel struct {
	modelsystem.Model
	transactionFn   func(context.Context, func(*gorm.DB) error) error
	getNodeConfigFn func(context.Context) ([]*modelsystem.System, error)
}

func (f fakeSystemModel) Transaction(ctx context.Context, fn func(*gorm.DB) error) error {
	if f.transactionFn == nil {
		panic("unexpected Transaction call")
	}
	return f.transactionFn(ctx, fn)
}

func (f fakeSystemModel) GetNodeConfig(ctx context.Context) ([]*modelsystem.System, error) {
	if f.getNodeConfigFn == nil {
		panic("unexpected GetNodeConfig call")
	}
	return f.getNodeConfigFn(ctx)
}

func requireSystemCodeError(t *testing.T, err error, want uint32) {
	t.Helper()

	require.Error(t, err)

	var codeErr *xerr.CodeError
	require.ErrorAs(t, err, &codeErr)
	require.Equal(t, want, codeErr.GetErrCode())
}

func testAdminSystemDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := fmt.Sprintf(
		"file:%s?mode=memory&cache=shared",
		strings.NewReplacer("/", "_", " ", "_").Replace(t.Name()),
	)
	db, err := gorm.Open(sqliteDriver.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	require.NoError(t, err)
	return db
}

func markIdentitySystemRevisionApplied(t *testing.T, db *gorm.DB) {
	t.Helper()

	require.NoError(t, db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_registry (
			id text primary key,
			source text not null,
			state text not null,
			checksum text,
			applied_at datetime not null,
			created_at datetime,
			updated_at datetime
		)
	`).Error)

	now := time.Now()
	require.NoError(t, db.Exec(`
		INSERT INTO schema_registry (id, source, state, checksum, applied_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, "0002_identity_system", "embedded", "applied", "", now, now, now).Error)
}
