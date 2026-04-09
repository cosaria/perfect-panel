package identity

import (
	"fmt"
	"strings"
	"testing"
	"time"

	sqliteDriver "gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func TestRepositoryAvailableRequiresAppliedRevision(t *testing.T) {
	t.Parallel()

	db := openIdentityTestDB(t)
	if err := db.AutoMigrate(&User{}, &AuthIdentity{}, &UserDevice{}); err != nil {
		t.Fatalf("migrate identity tables: %v", err)
	}

	repo := NewRepository(db)
	if repo.Available() {
		t.Fatalf("expected repository to stay disabled before revision is applied")
	}

	markIdentitySystemRevision(t, db, revisionStateApplied)
	if !repo.Available() {
		t.Fatalf("expected repository to become available after applied revision marker")
	}
}

func TestRepositoryAvailableRejectsPendingRevision(t *testing.T) {
	t.Parallel()

	db := openIdentityTestDB(t)
	if err := db.AutoMigrate(&User{}, &AuthIdentity{}, &UserDevice{}); err != nil {
		t.Fatalf("migrate identity tables: %v", err)
	}
	markIdentitySystemRevision(t, db, "pending")

	if NewRepository(db).Available() {
		t.Fatalf("expected pending revision marker to keep repository disabled")
	}
}

func openIdentityTestDB(t *testing.T) *gorm.DB {
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

func markIdentitySystemRevision(t *testing.T, db *gorm.DB, state string) {
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
