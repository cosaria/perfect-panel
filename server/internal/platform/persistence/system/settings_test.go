package system

import (
	"fmt"
	"strings"
	"testing"
	"time"

	sqliteDriver "gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func TestRepositorySchemaChecksRequireAppliedRevision(t *testing.T) {
	t.Parallel()

	db := openSystemTestDB(t)
	if err := db.AutoMigrate(&Setting{}, &VerificationPolicy{}, &AuthProvider{}, &AuthProviderConfig{}); err != nil {
		t.Fatalf("migrate system tables: %v", err)
	}

	repo := NewRepository(db)
	if repo.HasSettingsSchema() {
		t.Fatalf("expected settings schema to stay disabled before revision marker")
	}
	if repo.HasAuthProviderSchema() {
		t.Fatalf("expected auth provider schema to stay disabled before revision marker")
	}

	markSystemIdentityRevision(t, db, revisionStateApplied)
	if !repo.HasSettingsSchema() {
		t.Fatalf("expected settings schema after applied revision marker")
	}
	if !repo.HasAuthProviderSchema() {
		t.Fatalf("expected auth provider schema after applied revision marker")
	}
}

func TestRepositorySchemaChecksRejectPendingRevision(t *testing.T) {
	t.Parallel()

	db := openSystemTestDB(t)
	if err := db.AutoMigrate(&Setting{}, &VerificationPolicy{}, &AuthProvider{}, &AuthProviderConfig{}); err != nil {
		t.Fatalf("migrate system tables: %v", err)
	}
	markSystemIdentityRevision(t, db, "pending")

	repo := NewRepository(db)
	if repo.HasSettingsSchema() {
		t.Fatalf("expected pending revision marker to keep settings schema disabled")
	}
	if repo.HasAuthProviderSchema() {
		t.Fatalf("expected pending revision marker to keep auth provider schema disabled")
	}
}

func openSystemTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := fmt.Sprintf(
		"file:%s?mode=memory&cache=shared",
		strings.NewReplacer("/", "_", " ", "_").Replace(t.Name()),
	)
	db, err := gorm.Open(sqliteDriver.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		t.Fatalf("open sqlite db: %v", err)
	}
	return db
}

func markSystemIdentityRevision(t *testing.T, db *gorm.DB, state string) {
	t.Helper()

	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_registry (
			id text primary key,
			source text not null,
			state text not null,
			checksum text,
			applied_at datetime not null,
			created_at datetime,
			updated_at datetime
		)
	`).Error; err != nil {
		t.Fatalf("create schema registry table: %v", err)
	}

	now := time.Now()
	if err := db.Exec(`
		INSERT INTO schema_registry (id, source, state, checksum, applied_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, identitySystemRevisionName, "embedded", state, "", now, now, now).Error; err != nil {
		t.Fatalf("insert revision marker: %v", err)
	}
}
