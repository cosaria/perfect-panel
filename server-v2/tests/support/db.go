package support

import "testing"

// DBStub 是供测试装配使用的最小数据库占位。
type DBStub struct{}

// NewDBStub 返回数据库测试占位实例。
func NewDBStub(t *testing.T) *DBStub {
	t.Helper()
	return &DBStub{}
}
