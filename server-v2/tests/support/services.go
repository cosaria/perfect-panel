package support

import "testing"

// ServiceEnv 是供测试复用的最小服务环境。
type ServiceEnv struct {
	Runtime RuntimeEnv
	DB      *DBStub
}

// NewServiceEnv 创建可复用的最小服务环境。
func NewServiceEnv(t *testing.T) *ServiceEnv {
	t.Helper()

	rt := NewRuntimeEnv(t)
	return &ServiceEnv{
		Runtime: rt,
		DB:      NewDBStub(t, "stub://"+rt.TempDir),
	}
}

// WriteFixture 写入 JSON fixture 并返回路径。
func (e *ServiceEnv) WriteFixture(t *testing.T, fileName string, value any) string {
	t.Helper()
	return WriteJSONFixture(t, e.Runtime.TempDir, fileName, value)
}
