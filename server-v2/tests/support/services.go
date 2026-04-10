package support

import "testing"

// ServicesStub 是供后续测试复用的服务占位容器。
type ServicesStub struct {
	DB *DBStub
}

// NewServicesStub 创建最小测试服务集合。
func NewServicesStub(t *testing.T) *ServicesStub {
	t.Helper()
	return &ServicesStub{
		DB: NewDBStub(t),
	}
}
