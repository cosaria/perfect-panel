package configinit

import (
	"fmt"
	"strings"
	"testing"

	"github.com/perfect-panel/server/config"
	modelauth "github.com/perfect-panel/server/internal/platform/persistence/auth"
	modelclient "github.com/perfect-panel/server/internal/platform/persistence/client"
	modelidentity "github.com/perfect-panel/server/internal/platform/persistence/identity"
	"github.com/perfect-panel/server/internal/platform/persistence/schema"
	schemarevisions "github.com/perfect-panel/server/internal/platform/persistence/schema/revisions"
	modelsubscription "github.com/perfect-panel/server/internal/platform/persistence/subscription"
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
	schemarevisions.RegisterEmbedded()

	Migrate(Deps{DB: db, Config: cfg})
	Migrate(Deps{DB: db, Config: cfg})

	var registryCount int64
	if err := db.Model(&schema.Registry{}).Count(&registryCount).Error; err != nil {
		t.Fatalf("count schema registry rows: %v", err)
	}
	if registryCount != int64(len(schema.RegisteredRevisions())) {
		t.Fatalf("expected %d schema revision rows, got %d", len(schema.RegisteredRevisions()), registryCount)
	}

	var authProviderCount int64
	if err := db.Model(&modelauth.Auth{}).Count(&authProviderCount).Error; err != nil {
		t.Fatalf("count auth providers: %v", err)
	}
	if authProviderCount == 0 {
		t.Fatal("expected migrate to seed legacy auth methods")
	}

	var normalizedAuthProviderCount int64
	if err := db.Model(&modelsystem.AuthProvider{}).Count(&normalizedAuthProviderCount).Error; err != nil {
		t.Fatalf("count normalized auth providers: %v", err)
	}
	if normalizedAuthProviderCount == 0 {
		t.Fatal("expected migrate to seed normalized auth providers")
	}

	var systemCount int64
	if err := db.Model(&modelsystem.System{}).Count(&systemCount).Error; err != nil {
		t.Fatalf("count system rows: %v", err)
	}
	if systemCount == 0 {
		t.Fatal("expected migrate to seed legacy system rows")
	}

	var normalizedSystemCount int64
	if err := db.Model(&modelsystem.Setting{}).Count(&normalizedSystemCount).Error; err != nil {
		t.Fatalf("count normalized system rows: %v", err)
	}
	if normalizedSystemCount == 0 {
		t.Fatal("expected migrate to seed normalized system rows")
	}

	var verificationPolicyCount int64
	if err := db.Model(&modelsystem.VerificationPolicy{}).Count(&verificationPolicyCount).Error; err != nil {
		t.Fatalf("count verification policies: %v", err)
	}
	if verificationPolicyCount == 0 {
		t.Fatal("expected migrate to seed verification policies")
	}

	var userCount int64
	if err := db.Model(&modeluser.User{}).Count(&userCount).Error; err != nil {
		t.Fatalf("count admin users: %v", err)
	}
	if userCount != 1 {
		t.Fatalf("expected one admin user, got %d", userCount)
	}

	var identityUserCount int64
	if err := db.Model(&modelidentity.User{}).Count(&identityUserCount).Error; err != nil {
		t.Fatalf("count normalized admin users: %v", err)
	}
	if identityUserCount != 1 {
		t.Fatalf("expected one normalized admin user, got %d", identityUserCount)
	}

	var userAuthCount int64
	if err := db.Model(&modeluser.AuthMethods{}).Count(&userAuthCount).Error; err != nil {
		t.Fatalf("count admin auth methods: %v", err)
	}
	if userAuthCount != 1 {
		t.Fatalf("expected one admin auth method, got %d", userAuthCount)
	}

	var identityAuthCount int64
	if err := db.Model(&modelidentity.AuthIdentity{}).Count(&identityAuthCount).Error; err != nil {
		t.Fatalf("count normalized admin auth methods: %v", err)
	}
	if identityAuthCount != 1 {
		t.Fatalf("expected one normalized admin auth method, got %d", identityAuthCount)
	}

	if !db.Migrator().HasTable(&modelsubscription.Subscription{}) {
		t.Fatal("expected migrate to apply subscription revisions")
	}
	if !db.Migrator().HasTable(&modelclient.SubscribeApplication{}) {
		t.Fatal("expected migrate to provision subscribe applications table")
	}
}

func TestBootstrapInitDatabaseSeedsAdminIdempotently(t *testing.T) {
	db := testConfigInitDB(t)
	schemarevisions.RegisterEmbedded()
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
	if registryCount != int64(len(schema.RegisteredRevisions())) {
		t.Fatalf("expected %d schema revision rows, got %d", len(schema.RegisteredRevisions()), registryCount)
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

	if !db.Migrator().HasTable(&modelsubscription.Subscription{}) {
		t.Fatal("expected bootstrap init to apply subscription revisions")
	}
	if !db.Migrator().HasTable(&modelclient.SubscribeApplication{}) {
		t.Fatal("expected bootstrap init to provision subscribe applications table")
	}
}
