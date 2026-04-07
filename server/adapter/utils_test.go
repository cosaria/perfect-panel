package adapter

import (
	"os"
	"testing"

	"github.com/perfect-panel/server/models/node"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestAdapterProxy(t *testing.T) {
	nodes := getNodes(t)
	adapter := NewAdapter(tpl, WithServers(nodes), WithUserInfo(User{Password: "test-password"}))

	proxies, err := adapter.Proxies(nodes)
	if err != nil {
		t.Fatalf("failed to adapt nodes: %v", err)
	}
	if len(proxies) == 0 {
		t.Fatal("no proxies generated")
	}
	for _, proxy := range proxies {
		t.Logf("[测试] 适配节点 %s 成功: %+v", proxy.Name, proxy)
	}

	if _, err := adapter.Client(); err != nil {
		t.Fatalf("failed to build adapter client: %v", err)
	}

	if len(adapter.Servers) != len(nodes) {
		t.Fatalf("adapter server count mismatch: got %d want %d", len(adapter.Servers), len(nodes))
	}

	if adapter.UserInfo.Password == "" {
		t.Fatal("adapter user info was not applied")
	}

	if adapter.OutputFormat == "" {
		t.Fatal("adapter output format should not be empty")
	}

	if adapter.ClientTemplate == "" {
		t.Fatal("adapter client template should not be empty")
	}

	if adapter.Type != "" {
		t.Fatalf("unexpected adapter type: %q", adapter.Type)
	}

	if adapter.Params != nil {
		t.Fatalf("unexpected adapter params: %+v", adapter.Params)
	}

	if adapter.SiteName != "" {
		t.Fatalf("unexpected adapter site name: %q", adapter.SiteName)
	}

	if adapter.SubscribeName != "" {
		t.Fatalf("unexpected adapter subscribe name: %q", adapter.SubscribeName)
	}

}

func getNodes(t *testing.T) []*node.Node {
	t.Helper()

	dsn := os.Getenv("PPANEL_TEST_MYSQL_DSN")
	if dsn == "" {
		t.Skip("set PPANEL_TEST_MYSQL_DSN to run adapter integration tests")
	}

	db, err := connectMySQL(dsn)
	if err != nil {
		t.Fatalf("failed to connect mysql: %v", err)
	}
	var nodes []*node.Node
	if err = db.Preload("Server").Find(&nodes).Error; err != nil {
		t.Fatalf("failed to load nodes: %v", err)
	}
	if len(nodes) == 0 {
		t.Fatal("no nodes found")
	}
	return nodes
}

func connectMySQL(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: dsn,
	}), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
