package identity_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/perfect-panel/server/internal/platform/persistence/schema"
	schemarevisions "github.com/perfect-panel/server/internal/platform/persistence/schema/revisions"
	"github.com/perfect-panel/server/internal/platform/persistence/subscribe"
	"github.com/perfect-panel/server/internal/platform/persistence/user"
	"github.com/redis/go-redis/v9"
	sqliteDriver "gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func TestIdentityCompatibility(t *testing.T) {
	t.Parallel()

	db := testIdentityDB(t)
	rds := testIdentityRedis(t)
	ctx := context.Background()

	emailVerified := true
	adminEnabled := true
	onlyFirstPurchase := false

	legacyUser := &user.User{
		Password:          "encoded-password",
		ReferCode:         "REF-CODE-001",
		Enable:            &adminEnabled,
		IsAdmin:           &adminEnabled,
		OnlyFirstPurchase: &onlyFirstPurchase,
		Avatar:            "https://example.com/avatar.png",
	}
	if err := db.Create(legacyUser).Error; err != nil {
		t.Fatalf("create legacy user: %v", err)
	}

	legacyAuthMethods := []*user.AuthMethods{
		{
			UserId:         legacyUser.Id,
			AuthType:       "email",
			AuthIdentifier: "worker@example.com",
			Verified:       emailVerified,
		},
		{
			UserId:         legacyUser.Id,
			AuthType:       "google",
			AuthIdentifier: "google-open-id",
			Verified:       true,
		},
	}
	for _, item := range legacyAuthMethods {
		if err := db.Create(item).Error; err != nil {
			t.Fatalf("create legacy auth method: %v", err)
		}
	}

	legacyDevice := &user.Device{
		UserId:     legacyUser.Id,
		Ip:         "127.0.0.1",
		UserAgent:  "compat-test",
		Identifier: "device-001",
		Enabled:    true,
		Online:     false,
	}
	if err := db.Create(legacyDevice).Error; err != nil {
		t.Fatalf("create legacy device: %v", err)
	}

	legacyPlan := &subscribe.Subscribe{
		Id:       1,
		Name:     "Compat Plan",
		Language: "en",
		UnitTime: "month",
	}
	if err := db.Exec(`
		INSERT INTO subscribe (id, name, language, unit_time, created_at, updated_at)
		VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, legacyPlan.Id, legacyPlan.Name, legacyPlan.Language, legacyPlan.UnitTime).Error; err != nil {
		t.Fatalf("create legacy subscribe plan: %v", err)
	}

	legacyUserSubscribe := &user.Subscribe{
		Id:          1,
		UserId:      legacyUser.Id,
		OrderId:     1,
		SubscribeId: legacyPlan.Id,
		Token:       "token-compat-001",
		UUID:        "uuid-compat-001",
		Status:      1,
	}
	if err := db.Exec(`
		INSERT INTO user_subscribe (
			id, user_id, order_id, subscribe_id, start_time, traffic, download, upload, token, uuid, status, note, created_at, updated_at
		) VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP, 0, 0, 0, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, legacyUserSubscribe.Id, legacyUserSubscribe.UserId, legacyUserSubscribe.OrderId, legacyUserSubscribe.SubscribeId, legacyUserSubscribe.Token, legacyUserSubscribe.UUID, legacyUserSubscribe.Status, legacyUserSubscribe.Note).Error; err != nil {
		t.Fatalf("create legacy user subscribe: %v", err)
	}

	if err := schema.ApplyRevisions(db, schema.DefaultRevisionSource); err != nil {
		t.Fatalf("apply revisions: %v", err)
	}

	if err := db.Migrator().DropTable("user_device", "user_auth_methods", "user"); err != nil {
		t.Fatalf("drop legacy identity tables: %v", err)
	}

	model := user.NewModel(db, rds)

	t.Run("FindOneByEmail 保持语义", func(t *testing.T) {
		found, err := model.FindOneByEmail(ctx, "worker@example.com")
		if err != nil {
			t.Fatalf("FindOneByEmail returned error: %v", err)
		}
		if found.Id != legacyUser.Id {
			t.Fatalf("expected user id %d, got %d", legacyUser.Id, found.Id)
		}
		if found.ReferCode != legacyUser.ReferCode {
			t.Fatalf("expected refer code %q, got %q", legacyUser.ReferCode, found.ReferCode)
		}
		if len(found.AuthMethods) != 2 {
			t.Fatalf("expected 2 auth methods, got %d", len(found.AuthMethods))
		}
		if len(found.UserDevices) != 1 {
			t.Fatalf("expected 1 user device, got %d", len(found.UserDevices))
		}
	})

	t.Run("FindUserAuthMethods 与 InsertUserAuthMethods 保持语义", func(t *testing.T) {
		methods, err := model.FindUserAuthMethods(ctx, legacyUser.Id)
		if err != nil {
			t.Fatalf("FindUserAuthMethods returned error: %v", err)
		}
		if len(methods) != 2 {
			t.Fatalf("expected 2 auth methods before insert, got %d", len(methods))
		}

		inserted := &user.AuthMethods{
			UserId:         legacyUser.Id,
			AuthType:       "mobile",
			AuthIdentifier: "+8613012345678",
			Verified:       false,
		}
		if err := model.InsertUserAuthMethods(ctx, inserted); err != nil {
			t.Fatalf("InsertUserAuthMethods returned error: %v", err)
		}
		if inserted.Id == 0 {
			t.Fatalf("expected inserted auth method id to be populated")
		}

		methods, err = model.FindUserAuthMethods(ctx, legacyUser.Id)
		if err != nil {
			t.Fatalf("FindUserAuthMethods after insert returned error: %v", err)
		}
		if len(methods) != 3 {
			t.Fatalf("expected 3 auth methods after insert, got %d", len(methods))
		}
	})

	t.Run("QueryPageList 与统计接口保持可用", func(t *testing.T) {
		list, total, err := model.QueryPageList(ctx, 1, 10, &user.UserFilterParams{
			Search: "worker@example.com",
			Order:  "DESC",
		})
		if err != nil {
			t.Fatalf("QueryPageList returned error: %v", err)
		}
		if total != 1 {
			t.Fatalf("expected total 1, got %d", total)
		}
		if len(list) != 1 || list[0].Id != legacyUser.Id {
			t.Fatalf("expected one matching user %d, got %+v", legacyUser.Id, list)
		}

		count, err := model.QueryResisterUserTotal(ctx)
		if err != nil {
			t.Fatalf("QueryResisterUserTotal returned error: %v", err)
		}
		if count != 1 {
			t.Fatalf("expected register total 1, got %d", count)
		}
	})

	t.Run("BatchDeleteUser 在新 schema 下仍可工作", func(t *testing.T) {
		if err := model.BatchDeleteUser(ctx, []int64{legacyUser.Id}); err != nil {
			t.Fatalf("BatchDeleteUser returned error: %v", err)
		}

		count, err := model.QueryResisterUserTotal(ctx)
		if err != nil {
			t.Fatalf("QueryResisterUserTotal after delete returned error: %v", err)
		}
		if count != 0 {
			t.Fatalf("expected register total 0 after delete, got %d", count)
		}

		list, total, err := model.QueryPageList(ctx, 1, 10, &user.UserFilterParams{})
		if err != nil {
			t.Fatalf("QueryPageList after delete returned error: %v", err)
		}
		if total != 0 || len(list) != 0 {
			t.Fatalf("expected deleted user to disappear from default list, got total=%d len=%d", total, len(list))
		}
	})

	t.Run("FindOneDevice 保持语义", func(t *testing.T) {
		device, err := model.FindOneDevice(ctx, legacyDevice.Id)
		if err != nil {
			t.Fatalf("FindOneDevice returned error: %v", err)
		}
		if device.Identifier != legacyDevice.Identifier {
			t.Fatalf("expected identifier %q, got %q", legacyDevice.Identifier, device.Identifier)
		}
		if device.UserId != legacyUser.Id {
			t.Fatalf("expected user id %d, got %d", legacyUser.Id, device.UserId)
		}
	})

	t.Run("FindOneSubscribeDetailsById 保持语义", func(t *testing.T) {
		detail, err := model.FindOneSubscribeDetailsById(ctx, legacyUserSubscribe.Id)
		if err != nil {
			t.Fatalf("FindOneSubscribeDetailsById returned error: %v", err)
		}
		if detail.User == nil || detail.User.Id != legacyUser.Id {
			t.Fatalf("expected subscribe detail to preload legacy user id %d, got %+v", legacyUser.Id, detail.User)
		}
		if detail.Subscribe == nil || detail.Subscribe.Id != legacyPlan.Id {
			t.Fatalf("expected subscribe detail to preload subscribe id %d, got %+v", legacyPlan.Id, detail.Subscribe)
		}
	})
}

func testIdentityDB(t *testing.T) *gorm.DB {
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
	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS subscribe (
			id integer primary key,
			name text not null,
			language text not null,
			unit_time text not null,
			created_at datetime,
			updated_at datetime
		)
	`).Error; err != nil {
		t.Fatalf("migrate legacy subscribe table: %v", err)
	}
	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS user_subscribe (
			id integer primary key autoincrement,
			user_id integer not null,
			order_id integer not null,
			subscribe_id integer not null,
			start_time datetime,
			expire_time datetime,
			finished_at datetime,
			traffic integer not null default 0,
			download integer not null default 0,
			upload integer not null default 0,
			token text not null,
			uuid text not null,
			status integer not null default 0,
			note text not null default '',
			created_at datetime,
			updated_at datetime
		)
	`).Error; err != nil {
		t.Fatalf("migrate legacy user subscribe table: %v", err)
	}

	return db
}

func testIdentityRedis(t *testing.T) *redis.Client {
	t.Helper()

	mr := miniredis.RunT(t)
	return redis.NewClient(&redis.Options{Addr: mr.Addr()})
}
