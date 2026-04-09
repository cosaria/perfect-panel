package catalog_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/perfect-panel/server/internal/platform/persistence/node"
	"github.com/perfect-panel/server/internal/platform/persistence/schema"
	schemarevisions "github.com/perfect-panel/server/internal/platform/persistence/schema/revisions"
	"github.com/perfect-panel/server/internal/platform/persistence/subscribe"
	"github.com/redis/go-redis/v9"
	sqliteDriver "gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func TestFilterListUsesRelationsWhenLegacyColumnsAreBlank(t *testing.T) {
	t.Parallel()

	db := openCatalogTestDB(t)
	rds := openCatalogTestRedis(t)
	ctx := context.Background()

	enabled := true
	serverRow := &node.Server{
		Id:        1,
		Name:      "edge-1",
		Address:   "198.51.100.10",
		Country:   "US",
		City:      "San Jose",
		Protocols: `[{"type":"vless","port":443,"enable":true}]`,
	}
	if err := db.Create(serverRow).Error; err != nil {
		t.Fatalf("create server: %v", err)
	}

	nodeRows := []*node.Node{
		{
			Id:       1,
			Name:     "node-tagged",
			Tags:     "premium,video",
			Port:     443,
			Address:  "198.51.100.10",
			ServerId: serverRow.Id,
			Protocol: "vless",
			Enabled:  &enabled,
		},
		{
			Id:       2,
			Name:     "node-explicit",
			Tags:     "standard",
			Port:     8443,
			Address:  "198.51.100.11",
			ServerId: serverRow.Id,
			Protocol: "vless",
			Enabled:  &enabled,
		},
	}
	for _, item := range nodeRows {
		if err := db.Create(item).Error; err != nil {
			t.Fatalf("create node %d: %v", item.Id, err)
		}
	}

	subscribeRows := []*subscribe.Subscribe{
		{
			Id:       11,
			Name:     "Tag Plan",
			Language: "zh-CN",
			UnitTime: "month",
			NodeTags: "premium",
		},
		{
			Id:       12,
			Name:     "Explicit Plan",
			Language: "zh-CN",
			UnitTime: "month",
			Nodes:    "2",
		},
	}
	for _, item := range subscribeRows {
		if err := db.Create(item).Error; err != nil {
			t.Fatalf("create subscribe %d: %v", item.Id, err)
		}
	}

	if err := schema.ApplyRevisions(db, schema.DefaultRevisionSource); err != nil {
		t.Fatalf("apply revisions: %v", err)
	}

	if err := db.Model(&subscribe.Subscribe{}).
		Where("id IN ?", []int64{11, 12}).
		Updates(map[string]any{"nodes": "", "node_tags": ""}).Error; err != nil {
		t.Fatalf("blank legacy selector columns: %v", err)
	}

	model := subscribe.NewModel(db, rds)

	t.Run("节点 ID 关系在旧列清空后仍然可查", func(t *testing.T) {
		total, list, err := model.FilterList(ctx, &subscribe.FilterParams{
			Page: 1,
			Size: 10,
			Node: []int64{2},
		})
		if err != nil {
			t.Fatalf("FilterList by node returned error: %v", err)
		}
		if total != 1 {
			t.Fatalf("expected one explicit-node plan, got total=%d", total)
		}
		if len(list) != 1 || list[0].Id != 12 {
			t.Fatalf("expected explicit-node plan 12, got %+v", list)
		}
	})

	t.Run("标签关系在旧列清空后仍然可查", func(t *testing.T) {
		total, list, err := model.FilterList(ctx, &subscribe.FilterParams{
			Page: 1,
			Size: 10,
			Tags: []string{"premium"},
		})
		if err != nil {
			t.Fatalf("FilterList by tags returned error: %v", err)
		}
		if total != 1 {
			t.Fatalf("expected one tag-driven plan, got total=%d", total)
		}
		if len(list) != 1 || list[0].Id != 11 {
			t.Fatalf("expected tag-driven plan 11, got %+v", list)
		}
	})
}

func TestUpdateSubscribeClearsOldAndNewServerCachesAfterRelationSync(t *testing.T) {
	t.Parallel()

	db := openCatalogTestDB(t)
	rds := openCatalogTestRedis(t)
	ctx := context.Background()

	enabled := true
	serverRows := []*node.Server{
		{
			Id:        1,
			Name:      "edge-1",
			Address:   "198.51.100.21",
			Country:   "US",
			City:      "San Jose",
			Protocols: `[{"type":"vless","port":443,"enable":true}]`,
		},
		{
			Id:        2,
			Name:      "edge-2",
			Address:   "198.51.100.22",
			Country:   "US",
			City:      "Seattle",
			Protocols: `[{"type":"vless","port":443,"enable":true}]`,
		},
	}
	for _, item := range serverRows {
		if err := db.Create(item).Error; err != nil {
			t.Fatalf("create server %d: %v", item.Id, err)
		}
	}

	nodeRows := []*node.Node{
		{
			Id:       21,
			Name:     "node-old",
			Tags:     "starter",
			Port:     443,
			Address:  "198.51.100.21",
			ServerId: 1,
			Protocol: "vless",
			Enabled:  &enabled,
		},
		{
			Id:       22,
			Name:     "node-new",
			Tags:     "premium",
			Port:     8443,
			Address:  "198.51.100.22",
			ServerId: 2,
			Protocol: "vless",
			Enabled:  &enabled,
		},
	}
	for _, item := range nodeRows {
		if err := db.Create(item).Error; err != nil {
			t.Fatalf("create node %d: %v", item.Id, err)
		}
	}

	plan := &subscribe.Subscribe{
		Id:          31,
		Name:        "Cache Plan",
		Language:    "zh-CN",
		UnitTime:    "month",
		Nodes:       "21",
		SpeedLimit:  128,
		DeviceLimit: 3,
	}
	if err := db.Create(plan).Error; err != nil {
		t.Fatalf("create subscribe: %v", err)
	}

	if err := schema.ApplyRevisions(db, schema.DefaultRevisionSource); err != nil {
		t.Fatalf("apply revisions: %v", err)
	}

	model := subscribe.NewModel(db, rds)
	oldCacheKey := fmt.Sprintf("%s%d", node.ServerUserListCacheKey, 1)
	newCacheKey := fmt.Sprintf("%s%d", node.ServerUserListCacheKey, 2)
	if err := rds.Set(ctx, oldCacheKey, "stale-old", 0).Err(); err != nil {
		t.Fatalf("seed old cache key: %v", err)
	}
	if err := rds.Set(ctx, newCacheKey, "stale-new", 0).Err(); err != nil {
		t.Fatalf("seed new cache key: %v", err)
	}

	current, err := model.FindOne(ctx, plan.Id)
	if err != nil {
		t.Fatalf("FindOne returned error: %v", err)
	}
	current.Nodes = "22"
	current.NodeTags = ""
	if err := model.Update(ctx, current); err != nil {
		t.Fatalf("Update returned error: %v", err)
	}

	if rds.Exists(ctx, oldCacheKey).Val() != 0 {
		t.Fatalf("expected old server cache %q to be cleared", oldCacheKey)
	}
	if rds.Exists(ctx, newCacheKey).Val() != 0 {
		t.Fatalf("expected new server cache %q to be cleared", newCacheKey)
	}
}

func TestDeleteSubscribeClearsServerCacheFromExistingRelations(t *testing.T) {
	t.Parallel()

	db := openCatalogTestDB(t)
	rds := openCatalogTestRedis(t)
	ctx := context.Background()

	enabled := true
	serverRow := &node.Server{
		Id:        3,
		Name:      "edge-delete",
		Address:   "198.51.100.23",
		Country:   "US",
		City:      "Denver",
		Protocols: `[{"type":"vless","port":443,"enable":true}]`,
	}
	if err := db.Create(serverRow).Error; err != nil {
		t.Fatalf("create server: %v", err)
	}

	nodeRow := &node.Node{
		Id:       23,
		Name:     "node-delete",
		Tags:     "standard",
		Port:     443,
		Address:  "198.51.100.23",
		ServerId: serverRow.Id,
		Protocol: "vless",
		Enabled:  &enabled,
	}
	if err := db.Create(nodeRow).Error; err != nil {
		t.Fatalf("create node: %v", err)
	}

	plan := &subscribe.Subscribe{
		Id:       32,
		Name:     "Delete Plan",
		Language: "zh-CN",
		UnitTime: "month",
		Nodes:    "23",
	}
	if err := db.Create(plan).Error; err != nil {
		t.Fatalf("create subscribe: %v", err)
	}

	if err := schema.ApplyRevisions(db, schema.DefaultRevisionSource); err != nil {
		t.Fatalf("apply revisions: %v", err)
	}

	model := subscribe.NewModel(db, rds)
	cacheKey := fmt.Sprintf("%s%d", node.ServerUserListCacheKey, serverRow.Id)
	if err := rds.Set(ctx, cacheKey, "stale-delete", 0).Err(); err != nil {
		t.Fatalf("seed delete cache key: %v", err)
	}

	if err := model.Delete(ctx, plan.Id); err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}

	if rds.Exists(ctx, cacheKey).Val() != 0 {
		t.Fatalf("expected server cache %q to be cleared after delete", cacheKey)
	}
}

func openCatalogTestDB(t *testing.T) *gorm.DB {
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
	for _, stmt := range []string{
		`CREATE TABLE IF NOT EXISTS servers (
			id integer primary key,
			name text not null default '',
			country text not null default '',
			city text not null default '',
			address text not null default '',
			sort integer not null default 0,
			protocols text not null default '',
			last_reported_at datetime,
			created_at datetime default CURRENT_TIMESTAMP,
			updated_at datetime default CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS nodes (
			id integer primary key,
			name text not null default '',
			tags text not null default '',
			port integer not null default 0,
			address text not null default '',
			server_id integer not null default 0,
			protocol text not null default '',
			enabled numeric not null default 1,
			sort integer not null default 0,
			created_at datetime default CURRENT_TIMESTAMP,
			updated_at datetime default CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS subscribe (
			id integer primary key,
			name text not null default '',
			language text not null default '',
			description text default '',
			unit_price integer not null default 0,
			unit_time text not null default '',
			discount text default '',
			replacement integer not null default 0,
			inventory integer not null default -1,
			traffic integer not null default 0,
			speed_limit integer not null default 0,
			device_limit integer not null default 0,
			quota integer not null default 0,
			nodes text not null default '',
			node_tags text not null default '',
			show numeric not null default 0,
			sell numeric not null default 0,
			sort integer not null default 0,
			deduction_ratio integer not null default 0,
			allow_deduction numeric not null default 1,
			reset_cycle integer not null default 0,
			renewal_reset numeric not null default 0,
			show_original_price numeric not null default 1,
			created_at datetime default CURRENT_TIMESTAMP,
			updated_at datetime default CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS user_subscribe (
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
		)`,
	} {
		if err := db.Exec(stmt).Error; err != nil {
			t.Fatalf("create legacy test table failed: %v", err)
		}
	}

	return db
}

func openCatalogTestRedis(t *testing.T) *redis.Client {
	t.Helper()

	mr := miniredis.RunT(t)
	return redis.NewClient(&redis.Options{Addr: mr.Addr()})
}
