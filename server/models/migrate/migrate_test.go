package migrate

import (
	"os"
	"testing"

	"github.com/perfect-panel/server/models/node"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func testMySQLDSN(t *testing.T) string {
	t.Helper()

	dsn := os.Getenv("PPANEL_TEST_MYSQL_DSN")
	if dsn == "" {
		t.Skip("set PPANEL_TEST_MYSQL_DSN to run migrate integration tests")
	}
	return dsn
}

func TestMigrate(t *testing.T) {
	if os.Getenv("PPANEL_RUN_MIGRATE_TEST") != "1" {
		t.Skip("set PPANEL_RUN_MIGRATE_TEST=1 to run migration integration tests")
	}
	m := Migrate(testMySQLDSN(t))
	err := m.Migrate(2004)
	if err != nil {
		t.Errorf("failed to migrate: %v", err)
	} else {
		t.Log("migrate success")
	}
}
func TestMysql(t *testing.T) {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: testMySQLDSN(t),
	}))
	if err != nil {
		t.Fatalf("Failed to connect to MySQL: %v", err)
	}
	err = db.Migrator().AutoMigrate(&node.Node{})
	if err != nil {
		t.Fatalf("Failed to auto migrate: %v", err)
		return
	}
	t.Log("MySQL connection and migration successful")
}
