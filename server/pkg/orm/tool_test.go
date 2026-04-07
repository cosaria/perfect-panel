package orm

import (
	"os"
	"testing"

	"github.com/perfect-panel/server/models/task"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func testMySQLDSN(t *testing.T) string {
	t.Helper()

	dsn := os.Getenv("PPANEL_TEST_MYSQL_DSN")
	if dsn == "" {
		t.Skip("set PPANEL_TEST_MYSQL_DSN to run orm integration tests")
	}
	return dsn
}

func TestParseDSN(t *testing.T) {
	dsn := "root:mylove520@tcp(localhost:3306)/vpnboard"
	config := ParseDSN(dsn)
	if config == nil {
		t.Fatal("config is nil")
	}
	t.Log(config)
}

func TestPing(t *testing.T) {
	dsn := testMySQLDSN(t)
	status := Ping(dsn)
	if !status {
		t.Fatal("expected mysql ping to succeed")
	}
	t.Log(status)
}

func TestMysql(t *testing.T) {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: testMySQLDSN(t),
	}))
	if err != nil {
		t.Fatalf("Failed to connect to MySQL: %v", err)
	}
	err = db.Migrator().AutoMigrate(&task.Task{})
	if err != nil {
		t.Fatalf("Failed to auto migrate: %v", err)
		return
	}
	t.Log("MySQL connection and migration successful")
}
