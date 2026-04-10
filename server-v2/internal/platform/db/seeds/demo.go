package seeds

import (
	"context"
	"database/sql"
)

// ApplyDemo 预留演示数据入口，当前保持最小无副作用。
func ApplyDemo(_ context.Context, _ *sql.DB) error {
	return nil
}
