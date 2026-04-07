package tool

import (
	"os"
	"testing"
)

func testRedisURI(t *testing.T) string {
	t.Helper()

	uri := os.Getenv("PPANEL_TEST_REDIS_URI")
	if uri == "" {
		t.Skip("set PPANEL_TEST_REDIS_URI to run redis integration tests")
	}
	return uri
}

func TestParseRedisURI(t *testing.T) {
	uri := "redis://localhost:6379"
	addr, password, database, err := ParseRedisURI(uri)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(addr, password, database)
}

func TestRedisPing(t *testing.T) {
	uri := testRedisURI(t)
	addr, password, database, err := ParseRedisURI(uri)
	if err != nil {
		t.Fatal(err)
	}
	err = RedisPing(addr, password, database)
	if err != nil {
		t.Fatal(err)
	}
}
