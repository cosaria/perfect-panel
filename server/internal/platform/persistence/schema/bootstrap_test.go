package schema_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/perfect-panel/server/internal/platform/persistence/schema"
	schemarevisions "github.com/perfect-panel/server/internal/platform/persistence/schema/revisions"
	"github.com/perfect-panel/server/internal/platform/persistence/schema/seed"
	"github.com/perfect-panel/server/internal/platform/persistence/system"
	"github.com/perfect-panel/server/internal/platform/persistence/user"
	sqliteDriver "gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func testSchemaDB(t *testing.T) *gorm.DB {
	t.Helper()
	schemarevisions.RegisterEmbedded()
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

func TestBootstrapCreatesRegistryAndBaseline(t *testing.T) {
	db := testSchemaDB(t)
	if err := schema.Bootstrap(db, schema.DefaultRevisionSource); err != nil {
		t.Fatalf("bootstrap db: %v", err)
	}

	if !db.Migrator().HasTable(&schema.Registry{}) {
		t.Fatalf("expected schema registry table to exist")
	}
	if !db.Migrator().HasTable("user") {
		t.Fatalf("expected baseline to migrate user table")
	}
	if !db.Migrator().HasTable("user_auth_methods") {
		t.Fatalf("expected baseline to migrate user auth methods table")
	}

	var count int64
	if err := db.Model(&schema.Registry{}).Count(&count).Error; err != nil {
		t.Fatalf("count registry rows: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected one baseline revision, got %d", count)
	}
}

func TestBootstrapSupportsSeedWorkflow(t *testing.T) {
	db := testSchemaDB(t)
	if err := schema.Bootstrap(db, schema.DefaultRevisionSource); err != nil {
		t.Fatalf("bootstrap db: %v", err)
	}
	if err := seed.Site(db); err != nil {
		t.Fatalf("seed site: %v", err)
	}
	if err := seed.Site(db); err != nil {
		t.Fatalf("seed site second pass: %v", err)
	}
	if err := seed.Admin(db, "admin@ppanel.dev", "password"); err != nil {
		t.Fatalf("seed admin: %v", err)
	}
	if err := seed.Admin(db, "admin@ppanel.dev", "password"); err != nil {
		t.Fatalf("seed admin second pass: %v", err)
	}

	var systemCount int64
	if err := db.Model(&system.System{}).Count(&systemCount).Error; err != nil {
		t.Fatalf("count system rows: %v", err)
	}
	if systemCount == 0 {
		t.Fatal("expected site seed to create system rows")
	}

	var webAd system.System
	if err := db.Where("`key` = ?", "WebAD").First(&webAd).Error; err != nil {
		t.Fatalf("expected site seed to create WebAD system row: %v", err)
	}
	if webAd.Value != "false" {
		t.Fatalf("expected WebAD default false, got %q", webAd.Value)
	}

	var userCount int64
	if err := db.Model(&user.User{}).Count(&userCount).Error; err != nil {
		t.Fatalf("count user rows: %v", err)
	}
	if userCount != 1 {
		t.Fatalf("expected one admin user, got %d", userCount)
	}

	var authMethodCount int64
	if err := db.Model(&user.AuthMethods{}).Count(&authMethodCount).Error; err != nil {
		t.Fatalf("count auth method rows: %v", err)
	}
	if authMethodCount != 1 {
		t.Fatalf("expected one admin auth method, got %d", authMethodCount)
	}
}

func TestBootstrapRejectsUnknownRevisionSource(t *testing.T) {
	db := testSchemaDB(t)
	if err := schema.Bootstrap(db, "bogus"); err == nil {
		t.Fatal("expected bootstrap to reject unknown revision source")
	}
}
