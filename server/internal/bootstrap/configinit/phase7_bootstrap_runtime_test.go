package configinit

import (
	"fmt"
	"strings"
	"testing"

	"github.com/perfect-panel/server/config"
	modelauth "github.com/perfect-panel/server/internal/platform/persistence/auth"
	"github.com/perfect-panel/server/internal/platform/persistence/schema"
	modelsystem "github.com/perfect-panel/server/internal/platform/persistence/system"
	modeluser "github.com/perfect-panel/server/internal/platform/persistence/user"
	sqliteDriver "gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func testConfigInitDB(t *testing.T) *gorm.DB {
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

func TestMigrateBootstrapsAndSeedsIdempotently(t *testing.T) {
	db := testConfigInitDB(t)
	cfg := &config.Config{}
	cfg.Administrator.Email = "admin@ppanel.dev"
	cfg.Administrator.Password = "password"

	Migrate(Deps{DB: db, Config: cfg})
	Migrate(Deps{DB: db, Config: cfg})

	var registryCount int64
	if err := db.Model(&schema.Registry{}).Count(&registryCount).Error; err != nil {
		t.Fatalf("count schema registry rows: %v", err)
	}
	if registryCount != 1 {
		t.Fatalf("expected one schema revision row, got %d", registryCount)
	}

	var authProviderCount int64
	if err := db.Model(&modelauth.Auth{}).Count(&authProviderCount).Error; err != nil {
		t.Fatalf("count auth providers: %v", err)
	}
	if authProviderCount == 0 {
		t.Fatal("expected migrate to seed auth providers")
	}

	var systemCount int64
	if err := db.Model(&modelsystem.System{}).Count(&systemCount).Error; err != nil {
		t.Fatalf("count system rows: %v", err)
	}
	if systemCount == 0 {
		t.Fatal("expected migrate to seed system rows")
	}

	var userCount int64
	if err := db.Model(&modeluser.User{}).Count(&userCount).Error; err != nil {
		t.Fatalf("count admin users: %v", err)
	}
	if userCount != 1 {
		t.Fatalf("expected one admin user, got %d", userCount)
	}

	var userAuthCount int64
	if err := db.Model(&modeluser.AuthMethods{}).Count(&userAuthCount).Error; err != nil {
		t.Fatalf("count admin auth methods: %v", err)
	}
	if userAuthCount != 1 {
		t.Fatalf("expected one admin auth method, got %d", userAuthCount)
	}
}

func TestBootstrapInitDatabaseSeedsAdminIdempotently(t *testing.T) {
	db := testConfigInitDB(t)
	if err := bootstrapInitDatabase(db, "admin@ppanel.dev", "password"); err != nil {
		t.Fatalf("bootstrap init database first pass: %v", err)
	}
	if err := bootstrapInitDatabase(db, "admin@ppanel.dev", "password"); err != nil {
		t.Fatalf("bootstrap init database second pass: %v", err)
	}

	var registryCount int64
	if err := db.Model(&schema.Registry{}).Count(&registryCount).Error; err != nil {
		t.Fatalf("count schema registry rows: %v", err)
	}
	if registryCount != 1 {
		t.Fatalf("expected one schema revision row, got %d", registryCount)
	}

	var userCount int64
	if err := db.Model(&modeluser.User{}).Count(&userCount).Error; err != nil {
		t.Fatalf("count admin users: %v", err)
	}
	if userCount != 1 {
		t.Fatalf("expected one admin user, got %d", userCount)
	}

	var userAuthCount int64
	if err := db.Model(&modeluser.AuthMethods{}).Count(&userAuthCount).Error; err != nil {
		t.Fatalf("count admin auth methods: %v", err)
	}
	if userAuthCount != 1 {
		t.Fatalf("expected one admin auth method, got %d", userAuthCount)
	}
}
