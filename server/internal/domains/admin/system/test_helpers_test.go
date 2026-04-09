package system

import (
	"context"
	"testing"

	modelsystem "github.com/perfect-panel/server/internal/platform/persistence/system"
	"github.com/perfect-panel/server/internal/platform/support/xerr"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type fakeSystemModel struct {
	modelsystem.Model
	transactionFn   func(context.Context, func(*gorm.DB) error) error
	getNodeConfigFn func(context.Context) ([]*modelsystem.System, error)
}

func (f fakeSystemModel) Transaction(ctx context.Context, fn func(*gorm.DB) error) error {
	if f.transactionFn == nil {
		panic("unexpected Transaction call")
	}
	return f.transactionFn(ctx, fn)
}

func (f fakeSystemModel) GetNodeConfig(ctx context.Context) ([]*modelsystem.System, error) {
	if f.getNodeConfigFn == nil {
		panic("unexpected GetNodeConfig call")
	}
	return f.getNodeConfigFn(ctx)
}

func requireSystemCodeError(t *testing.T, err error, want uint32) {
	t.Helper()

	require.Error(t, err)

	var codeErr *xerr.CodeError
	require.ErrorAs(t, err, &codeErr)
	require.Equal(t, want, codeErr.GetErrCode())
}
