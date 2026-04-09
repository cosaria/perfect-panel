package portal

import (
	"context"
	"testing"

	serverconfig "github.com/perfect-panel/server/config"
	modelorder "github.com/perfect-panel/server/internal/platform/persistence/order"
	modelpayment "github.com/perfect-panel/server/internal/platform/persistence/payment"
	"github.com/perfect-panel/server/internal/platform/support/xerr"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type fakePortalOrderModel struct {
	modelorder.Model
	findOneByOrderNoFn  func(context.Context, string) (*modelorder.Order, error)
	updateFn            func(context.Context, *modelorder.Order, ...*gorm.DB) error
	updateOrderStatusFn func(context.Context, string, uint8, ...*gorm.DB) error
}

func (f fakePortalOrderModel) FindOneByOrderNo(ctx context.Context, orderNo string) (*modelorder.Order, error) {
	if f.findOneByOrderNoFn == nil {
		panic("unexpected FindOneByOrderNo call")
	}
	return f.findOneByOrderNoFn(ctx, orderNo)
}

func (f fakePortalOrderModel) Update(ctx context.Context, data *modelorder.Order, tx ...*gorm.DB) error {
	if f.updateFn == nil {
		panic("unexpected Update call")
	}
	return f.updateFn(ctx, data, tx...)
}

func (f fakePortalOrderModel) UpdateOrderStatus(ctx context.Context, orderNo string, status uint8, tx ...*gorm.DB) error {
	if f.updateOrderStatusFn == nil {
		panic("unexpected UpdateOrderStatus call")
	}
	return f.updateOrderStatusFn(ctx, orderNo, status, tx...)
}

type fakePortalPaymentModel struct {
	modelpayment.Model
	findOneFn func(context.Context, int64) (*modelpayment.Payment, error)
}

func (f fakePortalPaymentModel) FindOne(ctx context.Context, id int64) (*modelpayment.Payment, error) {
	if f.findOneFn == nil {
		panic("unexpected FindOne call")
	}
	return f.findOneFn(ctx, id)
}

func newPortalTestDeps() Deps {
	return Deps{
		Config: &serverconfig.Config{},
	}
}

func requirePortalCodeError(t *testing.T, err error, want uint32) {
	t.Helper()

	require.Error(t, err)

	var codeErr *xerr.CodeError
	require.ErrorAs(t, err, &codeErr)
	require.Equal(t, want, codeErr.GetErrCode())
}
