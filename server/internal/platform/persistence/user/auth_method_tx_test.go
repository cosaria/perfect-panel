package user_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/perfect-panel/server/internal/platform/persistence/schema"
	schemarevisions "github.com/perfect-panel/server/internal/platform/persistence/schema/revisions"
	"github.com/perfect-panel/server/internal/platform/persistence/user"
	"github.com/redis/go-redis/v9"
	sqliteDriver "gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func TestInsertUserAuthMethodsWorksWithinTransactionAfterUserInsert(t *testing.T) {
	t.Parallel()

	db := openUserTxTestDB(t)
	rds := miniredis.RunT(t)
	model := user.NewModel(db, redis.NewClient(&redis.Options{Addr: rds.Addr()}))
	ctx := context.Background()

	err := model.Transaction(ctx, func(tx *gorm.DB) error {
		enabled := true
		onlyFirstPurchase := false
		userInfo := &user.User{
			Enable:            &enabled,
			IsAdmin:           &enabled,
			OnlyFirstPurchase: &onlyFirstPurchase,
		}
		if err := model.Insert(ctx, userInfo, tx); err != nil {
			return err
		}
		authInfo := &user.AuthMethods{
			UserId:         userInfo.Id,
			AuthType:       "email",
			AuthIdentifier: "tx-user@example.com",
			Verified:       true,
		}
		if err := model.InsertUserAuthMethods(ctx, authInfo, tx); err != nil {
			return err
		}
		if authInfo.Id == 0 {
			t.Fatalf("expected auth method id to be populated inside transaction")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("transaction should succeed: %v", err)
	}

	found, err := model.FindOneByEmail(ctx, "tx-user@example.com")
	if err != nil {
		t.Fatalf("FindOneByEmail returned error: %v", err)
	}
	if len(found.AuthMethods) != 1 {
		t.Fatalf("expected one auth method after commit, got %d", len(found.AuthMethods))
	}
}

func TestInsertUserAuthMethodsDefersCacheInvalidationUntilCommit(t *testing.T) {
	t.Parallel()

	db := openUserTxTestDB(t)
	rds := miniredis.RunT(t)
	redisClient := redis.NewClient(&redis.Options{Addr: rds.Addr()})
	model := user.NewModel(db, redisClient)
	ctx := context.Background()

	enabled := true
	onlyFirstPurchase := false
	userInfo := &user.User{
		Enable:            &enabled,
		IsAdmin:           &enabled,
		OnlyFirstPurchase: &onlyFirstPurchase,
	}
	if err := model.Insert(ctx, userInfo); err != nil {
		t.Fatalf("insert committed user: %v", err)
	}

	userCacheKey := fmt.Sprintf("cache:user:id:%d", userInfo.Id)
	emailCacheKey := "cache:user:email:cached@example.com"
	rds.Set(userCacheKey, "stale-user")
	rds.Set(emailCacheKey, "stale-email")

	err := model.Transaction(ctx, func(tx *gorm.DB) error {
		authInfo := &user.AuthMethods{
			UserId:         userInfo.Id,
			AuthType:       "email",
			AuthIdentifier: "cached@example.com",
			Verified:       true,
		}
		if err := model.InsertUserAuthMethods(ctx, authInfo, tx); err != nil {
			return err
		}

		if !rds.Exists(userCacheKey) || !rds.Exists(emailCacheKey) {
			t.Fatalf("expected cache keys to survive until outer transaction commits")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("transaction should succeed: %v", err)
	}

	if rds.Exists(userCacheKey) || rds.Exists(emailCacheKey) {
		t.Fatalf("expected cache keys to be evicted after commit")
	}
}

func TestInsertDeviceDefersCacheInvalidationUntilCommit(t *testing.T) {
	t.Parallel()

	db := openUserTxTestDB(t)
	rds := miniredis.RunT(t)
	redisClient := redis.NewClient(&redis.Options{Addr: rds.Addr()})
	model := user.NewModel(db, redisClient)
	ctx := context.Background()

	enabled := true
	onlyFirstPurchase := false
	userInfo := &user.User{
		Enable:            &enabled,
		IsAdmin:           &enabled,
		OnlyFirstPurchase: &onlyFirstPurchase,
	}
	if err := model.Insert(ctx, userInfo); err != nil {
		t.Fatalf("insert committed user: %v", err)
	}

	deviceCacheKey := "cache:user:device:number:device-001"
	rds.Set(deviceCacheKey, "stale-device")

	err := model.Transaction(ctx, func(tx *gorm.DB) error {
		deviceInfo := &user.Device{
			UserId:     userInfo.Id,
			Ip:         "127.0.0.1",
			UserAgent:  "tx-test",
			Identifier: "device-001",
			Enabled:    true,
			Online:     false,
		}
		if err := model.InsertDevice(ctx, deviceInfo, tx); err != nil {
			return err
		}

		if !rds.Exists(deviceCacheKey) {
			t.Fatalf("expected device cache to survive until outer transaction commits")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("transaction should succeed: %v", err)
	}

	if rds.Exists(deviceCacheKey) {
		t.Fatalf("expected device cache to be evicted after commit")
	}
}

func openUserTxTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	schemarevisions.RegisterEmbedded()
	db, err := gorm.Open(sqliteDriver.Open("file:"+strings.ReplaceAll(t.Name(), "/", "_")+"?mode=memory&cache=shared"), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		t.Fatalf("open sqlite db: %v", err)
	}
	if err := schema.Bootstrap(db, schema.DefaultRevisionSource); err != nil {
		t.Fatalf("bootstrap schema: %v", err)
	}
	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS user_device (
			id integer primary key autoincrement,
			ip text not null,
			user_id integer not null,
			user_agent text,
			identifier text not null unique,
			online numeric not null default 0,
			enabled numeric not null default 1,
			created_at datetime,
			updated_at datetime
		)
	`).Error; err != nil {
		t.Fatalf("migrate legacy device table: %v", err)
	}
	if err := schema.ApplyRevisions(db, schema.DefaultRevisionSource); err != nil {
		t.Fatalf("apply revisions: %v", err)
	}
	return db
}
