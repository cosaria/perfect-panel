package user

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	sqliteDriver "gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func TestInsertSubscribeRollsBackWhenAssignmentSyncFails(t *testing.T) {
	t.Parallel()

	model, db := openUserSubscribeTxTestModel(t)
	model.assignmentSyncer = assignmentSyncerStub{syncErr: errors.New("sync failed")}
	ctx := context.Background()

	sub := &Subscribe{
		Id:          1,
		UserId:      11,
		OrderId:     21,
		SubscribeId: 31,
		Token:       "rollback-insert",
		UUID:        "rollback-insert",
		Status:      1,
	}
	if err := model.InsertSubscribe(ctx, sub); err == nil {
		t.Fatalf("expected InsertSubscribe to fail when assignment sync fails")
	}

	assertUserSubscribeCount(t, db, sub.Id, 0)
}

func TestUpdateSubscribeRollsBackWhenAssignmentSyncFails(t *testing.T) {
	t.Parallel()

	model, db := openUserSubscribeTxTestModel(t)
	ctx := context.Background()

	original := &Subscribe{
		Id:          2,
		UserId:      12,
		OrderId:     22,
		SubscribeId: 32,
		Token:       "rollback-update-old",
		UUID:        "rollback-update-old",
		Status:      1,
	}
	if err := db.Create(original).Error; err != nil {
		t.Fatalf("seed subscribe: %v", err)
	}

	model.assignmentSyncer = assignmentSyncerStub{syncErr: errors.New("sync failed")}
	updated := *original
	updated.Token = "rollback-update-new"
	if err := model.UpdateSubscribe(ctx, &updated); err == nil {
		t.Fatalf("expected UpdateSubscribe to fail when assignment sync fails")
	}

	var stored Subscribe
	if err := db.Where("id = ?", original.Id).First(&stored).Error; err != nil {
		t.Fatalf("reload subscribe: %v", err)
	}
	if stored.Token != original.Token {
		t.Fatalf("expected token to remain %q after rollback, got %q", original.Token, stored.Token)
	}
}

func TestDeleteSubscribeByIDRollsBackWhenAssignmentCleanupFails(t *testing.T) {
	t.Parallel()

	model, db := openUserSubscribeTxTestModel(t)
	ctx := context.Background()

	row := &Subscribe{
		Id:          3,
		UserId:      13,
		OrderId:     23,
		SubscribeId: 33,
		Token:       "rollback-delete",
		UUID:        "rollback-delete",
		Status:      1,
	}
	if err := db.Create(row).Error; err != nil {
		t.Fatalf("seed subscribe: %v", err)
	}

	model.assignmentSyncer = assignmentSyncerStub{deleteErr: errors.New("delete failed")}
	if err := model.DeleteSubscribeById(ctx, row.Id); err == nil {
		t.Fatalf("expected DeleteSubscribeById to fail when assignment cleanup fails")
	}

	assertUserSubscribeCount(t, db, row.Id, 1)
}

type assignmentSyncerStub struct {
	syncErr   error
	deleteErr error
}

func (s assignmentSyncerStub) SyncUserSubscription(_ context.Context, _, _ int64, _ uint8, _ ...*gorm.DB) error {
	return s.syncErr
}

func (s assignmentSyncerStub) DeleteUserSubscription(_ context.Context, _ int64, _ ...*gorm.DB) error {
	return s.deleteErr
}

func openUserSubscribeTxTestModel(t *testing.T) (*customUserModel, *gorm.DB) {
	t.Helper()

	db, err := gorm.Open(sqliteDriver.Open("file:"+strings.ReplaceAll(t.Name(), "/", "_")+"?mode=memory&cache=shared"), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		t.Fatalf("open sqlite db: %v", err)
	}
	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS user_subscribe (
			id integer primary key,
			user_id integer not null,
			order_id integer not null,
			subscribe_id integer not null,
			start_time datetime default CURRENT_TIMESTAMP,
			expire_time datetime,
			finished_at datetime,
			traffic integer not null default 0,
			download integer not null default 0,
			upload integer not null default 0,
			token text not null default '',
			uuid text not null default '',
			status integer not null default 0,
			note text not null default '',
			created_at datetime default CURRENT_TIMESTAMP,
			updated_at datetime default CURRENT_TIMESTAMP
		)
	`).Error; err != nil {
		t.Fatalf("create user_subscribe table: %v", err)
	}

	rds := miniredis.RunT(t)
	model, ok := NewModel(db, redis.NewClient(&redis.Options{Addr: rds.Addr()})).(*customUserModel)
	if !ok {
		t.Fatalf("expected *customUserModel from NewModel")
	}
	return model, db
}

func assertUserSubscribeCount(t *testing.T, db *gorm.DB, id int64, want int64) {
	t.Helper()

	var count int64
	if err := db.Model(&Subscribe{}).Where("id = ?", id).Count(&count).Error; err != nil {
		t.Fatalf("count subscribe rows: %v", err)
	}
	if count != want {
		t.Fatalf("expected subscribe row count %d for id %d, got %d", want, id, count)
	}
}
