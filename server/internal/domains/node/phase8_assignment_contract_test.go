package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	nodeModel "github.com/perfect-panel/server/internal/platform/persistence/node"
	"github.com/perfect-panel/server/internal/platform/persistence/schema"
	schemarevisions "github.com/perfect-panel/server/internal/platform/persistence/schema/revisions"
	subscribeModel "github.com/perfect-panel/server/internal/platform/persistence/subscribe"
	userModel "github.com/perfect-panel/server/internal/platform/persistence/user"
	"github.com/redis/go-redis/v9"
	sqliteDriver "gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func TestPhase8GetServerUserListUsesAssignmentsWhenLegacySelectorsAreBlank(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := openPhase8NodeDB(t)
	rds := openPhase8NodeRedis(t)

	enabled := true
	userRow := &userModel.User{
		Id:       7,
		Password: "hashed",
		Enable:   &enabled,
		IsAdmin:  boolPtr(false),
	}
	if err := db.Create(userRow).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}

	serverRow := &nodeModel.Server{
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

	nodeRow := &nodeModel.Node{
		Id:       2,
		Name:     "node-explicit",
		Tags:     "premium",
		Port:     443,
		Address:  "198.51.100.10",
		ServerId: serverRow.Id,
		Protocol: "vless",
		Enabled:  &enabled,
	}
	if err := db.Create(nodeRow).Error; err != nil {
		t.Fatalf("create node: %v", err)
	}

	plan := &subscribeModel.Subscribe{
		Id:          9,
		Name:        "Starter",
		Language:    "zh-CN",
		UnitTime:    "month",
		Nodes:       "2",
		SpeedLimit:  128,
		DeviceLimit: 3,
	}
	if err := db.Create(plan).Error; err != nil {
		t.Fatalf("create subscribe plan: %v", err)
	}

	userSubscribe := &userModel.Subscribe{
		Id:          88,
		UserId:      userRow.Id,
		OrderId:     1001,
		SubscribeId: plan.Id,
		Token:       "token-phase8",
		UUID:        "uuid-phase8",
		Status:      1,
	}
	if err := db.Create(userSubscribe).Error; err != nil {
		t.Fatalf("create user subscribe: %v", err)
	}

	if err := schema.ApplyRevisions(db, schema.DefaultRevisionSource); err != nil {
		t.Fatalf("apply revisions: %v", err)
	}

	if err := db.Model(&subscribeModel.Subscribe{}).
		Where("id = ?", plan.Id).
		Updates(map[string]any{"nodes": "", "node_tags": ""}).Error; err != nil {
		t.Fatalf("blank legacy selector columns: %v", err)
	}

	router := gin.New()
	router.GET("/node/users", GetServerUserListHandler(Deps{
		Redis:          rds,
		NodeModel:      nodeModel.NewModel(db, rds),
		SubscribeModel: subscribeModel.NewModel(db, rds),
		UserModel:      userModel.NewModel(db, rds),
	}))

	req := httptest.NewRequest(http.MethodGet, "/node/users?server_id=1&protocol=vless", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", resp.Code, resp.Body.String())
	}

	var body map[string][]map[string]any
	if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	users := body["users"]
	if len(users) != 1 {
		t.Fatalf("expected one assigned user, got %+v", users)
	}
	if id, ok := users[0]["id"].(float64); !ok || int64(id) != userSubscribe.Id {
		t.Fatalf("expected user subscribe id %d, got %+v", userSubscribe.Id, users[0]["id"])
	}
	if uuid, ok := users[0]["uuid"].(string); !ok || uuid != userSubscribe.UUID {
		t.Fatalf("expected uuid %q, got %+v", userSubscribe.UUID, users[0]["uuid"])
	}
	if speed, ok := users[0]["speed_limit"].(float64); !ok || int64(speed) != plan.SpeedLimit {
		t.Fatalf("expected speed limit %d, got %+v", plan.SpeedLimit, users[0]["speed_limit"])
	}
	if deviceLimit, ok := users[0]["device_limit"].(float64); !ok || int64(deviceLimit) != plan.DeviceLimit {
		t.Fatalf("expected device limit %d, got %+v", plan.DeviceLimit, users[0]["device_limit"])
	}

}

func TestPhase8GetServerUserListPromotesPendingAssignmentsToActive(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := openPhase8NodeDB(t)
	rds := openPhase8NodeRedis(t)

	enabled := true
	userRow := &userModel.User{
		Id:       8,
		Password: "hashed",
		Enable:   &enabled,
		IsAdmin:  boolPtr(false),
	}
	if err := db.Create(userRow).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}

	serverRow := &nodeModel.Server{
		Id:        1,
		Name:      "edge-pending",
		Address:   "198.51.100.15",
		Country:   "US",
		City:      "Portland",
		Protocols: `[{"type":"vless","port":443,"enable":true}]`,
	}
	if err := db.Create(serverRow).Error; err != nil {
		t.Fatalf("create server: %v", err)
	}

	nodeRow := &nodeModel.Node{
		Id:       5,
		Name:     "node-pending",
		Tags:     "starter",
		Port:     443,
		Address:  "198.51.100.15",
		ServerId: serverRow.Id,
		Protocol: "vless",
		Enabled:  &enabled,
	}
	if err := db.Create(nodeRow).Error; err != nil {
		t.Fatalf("create node: %v", err)
	}

	plan := &subscribeModel.Subscribe{
		Id:          13,
		Name:        "Pending Plan",
		Language:    "zh-CN",
		UnitTime:    "month",
		Nodes:       "5",
		SpeedLimit:  256,
		DeviceLimit: 4,
	}
	if err := db.Create(plan).Error; err != nil {
		t.Fatalf("create subscribe plan: %v", err)
	}

	userSubscribe := &userModel.Subscribe{
		Id:          90,
		UserId:      userRow.Id,
		OrderId:     1003,
		SubscribeId: plan.Id,
		Token:       "token-phase8-pending",
		UUID:        "uuid-phase8-pending",
		Status:      0,
	}
	if err := db.Create(userSubscribe).Error; err != nil {
		t.Fatalf("create pending user subscribe: %v", err)
	}

	if err := schema.ApplyRevisions(db, schema.DefaultRevisionSource); err != nil {
		t.Fatalf("apply revisions: %v", err)
	}
	if err := db.Model(&subscribeModel.Subscribe{}).
		Where("id = ?", plan.Id).
		Updates(map[string]any{"nodes": "", "node_tags": ""}).Error; err != nil {
		t.Fatalf("blank legacy selector columns: %v", err)
	}

	router := gin.New()
	router.GET("/node/users", GetServerUserListHandler(Deps{
		Redis:          rds,
		NodeModel:      nodeModel.NewModel(db, rds),
		SubscribeModel: subscribeModel.NewModel(db, rds),
		UserModel:      userModel.NewModel(db, rds),
	}))

	req := httptest.NewRequest(http.MethodGet, "/node/users?server_id=1&protocol=vless", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", resp.Code, resp.Body.String())
	}

	var body map[string][]map[string]any
	if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	users := body["users"]
	if len(users) != 1 {
		t.Fatalf("expected one user after pending activation, got %+v", users)
	}

	var status int64
	if err := db.Model(&userModel.Subscribe{}).
		Select("status").
		Where("id = ?", userSubscribe.Id).
		Scan(&status).Error; err != nil {
		t.Fatalf("query updated status: %v", err)
	}
	if status != 1 {
		t.Fatalf("expected pending subscription to be promoted to active, got %d", status)
	}
}

func TestPhase8GetServerUserListReturnsEmptyWhenNoAssignments(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := openPhase8NodeDB(t)
	rds := openPhase8NodeRedis(t)

	enabled := true
	serverRow := &nodeModel.Server{
		Id:        1,
		Name:      "edge-empty",
		Address:   "198.51.100.20",
		Country:   "US",
		City:      "Los Angeles",
		Protocols: `[{"type":"vless","port":443,"enable":true}]`,
	}
	if err := db.Create(serverRow).Error; err != nil {
		t.Fatalf("create server: %v", err)
	}
	nodeRow := &nodeModel.Node{
		Id:       3,
		Name:     "node-empty",
		Tags:     "empty",
		Port:     443,
		Address:  "198.51.100.20",
		ServerId: serverRow.Id,
		Protocol: "vless",
		Enabled:  &enabled,
	}
	if err := db.Create(nodeRow).Error; err != nil {
		t.Fatalf("create node: %v", err)
	}
	if err := schema.ApplyRevisions(db, schema.DefaultRevisionSource); err != nil {
		t.Fatalf("apply revisions: %v", err)
	}

	router := gin.New()
	router.GET("/node/users", GetServerUserListHandler(Deps{
		Redis:          rds,
		NodeModel:      nodeModel.NewModel(db, rds),
		SubscribeModel: subscribeModel.NewModel(db, rds),
		UserModel:      userModel.NewModel(db, rds),
	}))

	req := httptest.NewRequest(http.MethodGet, "/node/users?server_id=1&protocol=vless", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", resp.Code, resp.Body.String())
	}

	var body map[string][]map[string]any
	if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	users := body["users"]
	if len(users) != 0 {
		t.Fatalf("expected no users instead of pseudo fallback, got %+v", users)
	}
}

func TestPhase8UpdateNodeRemovesAssignmentsAfterTagShrink(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := openPhase8NodeDB(t)
	rds := openPhase8NodeRedis(t)

	enabled := true
	userRow := &userModel.User{
		Id:       9,
		Password: "hashed",
		Enable:   &enabled,
		IsAdmin:  boolPtr(false),
	}
	if err := db.Create(userRow).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}

	serverRow := &nodeModel.Server{
		Id:        1,
		Name:      "edge-tag-shrink",
		Address:   "198.51.100.30",
		Country:   "US",
		City:      "Seattle",
		Protocols: `[{"type":"vless","port":443,"enable":true}]`,
	}
	if err := db.Create(serverRow).Error; err != nil {
		t.Fatalf("create server: %v", err)
	}

	nodeRow := &nodeModel.Node{
		Id:       4,
		Name:     "node-tagged",
		Tags:     "premium,video",
		Port:     443,
		Address:  "198.51.100.30",
		ServerId: serverRow.Id,
		Protocol: "vless",
		Enabled:  &enabled,
	}
	if err := db.Create(nodeRow).Error; err != nil {
		t.Fatalf("create node: %v", err)
	}

	plan := &subscribeModel.Subscribe{
		Id:          10,
		Name:        "Tag Plan",
		Language:    "zh-CN",
		UnitTime:    "month",
		NodeTags:    "premium",
		SpeedLimit:  64,
		DeviceLimit: 2,
	}
	if err := db.Create(plan).Error; err != nil {
		t.Fatalf("create subscribe plan: %v", err)
	}

	userSubscribe := &userModel.Subscribe{
		Id:          89,
		UserId:      userRow.Id,
		OrderId:     1002,
		SubscribeId: plan.Id,
		Token:       "token-phase8-shrink",
		UUID:        "uuid-phase8-shrink",
		Status:      1,
	}
	if err := db.Create(userSubscribe).Error; err != nil {
		t.Fatalf("create user subscribe: %v", err)
	}

	if err := schema.ApplyRevisions(db, schema.DefaultRevisionSource); err != nil {
		t.Fatalf("apply revisions: %v", err)
	}
	if err := db.Model(&subscribeModel.Subscribe{}).
		Where("id = ?", plan.Id).
		Updates(map[string]any{"nodes": "", "node_tags": ""}).Error; err != nil {
		t.Fatalf("blank legacy selector columns: %v", err)
	}

	nodeStore := nodeModel.NewModel(db, rds)
	nodeRow.Tags = "standard"
	if err := nodeStore.UpdateNode(t.Context(), nodeRow); err != nil {
		t.Fatalf("update node tags: %v", err)
	}

	router := gin.New()
	router.GET("/node/users", GetServerUserListHandler(Deps{
		Redis:          rds,
		NodeModel:      nodeStore,
		SubscribeModel: subscribeModel.NewModel(db, rds),
		UserModel:      userModel.NewModel(db, rds),
	}))

	req := httptest.NewRequest(http.MethodGet, "/node/users?server_id=1&protocol=vless", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", resp.Code, resp.Body.String())
	}

	var body map[string][]map[string]any
	if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if users := body["users"]; len(users) != 0 {
		t.Fatalf("expected assignments to be removed after tag shrink, got %+v", users)
	}
}

func openPhase8NodeDB(t *testing.T) *gorm.DB {
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

func openPhase8NodeRedis(t *testing.T) *redis.Client {
	t.Helper()

	mr := miniredis.RunT(t)
	return redis.NewClient(&redis.Options{Addr: mr.Addr()})
}

func boolPtr(v bool) *bool {
	return &v
}
