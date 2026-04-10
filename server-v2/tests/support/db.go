package support

import "testing"

// DBStub 是供测试装配使用的最小数据库占位。
type DBStub struct {
	DSN   string
	ready bool
}

// NewDBStub 返回数据库测试占位实例。
func NewDBStub(t *testing.T, dsn string) *DBStub {
	t.Helper()
	if dsn == "" {
		dsn = "stub://memory"
	}
	return &DBStub{
		DSN:   dsn,
		ready: true,
	}
}

// IsReady 返回 DB stub 的就绪状态。
func (d *DBStub) IsReady() bool {
	return d != nil && d.ready
}
