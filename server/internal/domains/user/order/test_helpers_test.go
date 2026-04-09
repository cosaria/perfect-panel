package order

import (
	"context"
	"testing"

	serverconfig "github.com/perfect-panel/server/config"
	modelcoupon "github.com/perfect-panel/server/models/coupon"
	modelpayment "github.com/perfect-panel/server/models/payment"
	modelsubscribe "github.com/perfect-panel/server/models/subscribe"
	modeluser "github.com/perfect-panel/server/models/user"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/stretchr/testify/require"
)

type fakeUserModel struct {
	modeluser.Model
	queryUserSubscribeFn func(context.Context, int64, ...int64) ([]*modeluser.SubscribeDetails, error)
}

func (f fakeUserModel) QueryUserSubscribe(ctx context.Context, userID int64, status ...int64) ([]*modeluser.SubscribeDetails, error) {
	if f.queryUserSubscribeFn == nil {
		panic("unexpected QueryUserSubscribe call")
	}
	return f.queryUserSubscribeFn(ctx, userID, status...)
}

type fakeSubscribeModel struct {
	modelsubscribe.Model
	findOneFn func(context.Context, int64) (*modelsubscribe.Subscribe, error)
}

func (f fakeSubscribeModel) FindOne(ctx context.Context, id int64) (*modelsubscribe.Subscribe, error) {
	if f.findOneFn == nil {
		panic("unexpected FindOne call")
	}
	return f.findOneFn(ctx, id)
}

type fakeCouponModel struct {
	modelcoupon.Model
	findOneByCodeFn func(context.Context, string) (*modelcoupon.Coupon, error)
}

func (f fakeCouponModel) FindOneByCode(ctx context.Context, code string) (*modelcoupon.Coupon, error) {
	if f.findOneByCodeFn == nil {
		panic("unexpected FindOneByCode call")
	}
	return f.findOneByCodeFn(ctx, code)
}

type fakePaymentModel struct {
	modelpayment.Model
	findOneFn func(context.Context, int64) (*modelpayment.Payment, error)
}

func (f fakePaymentModel) FindOne(ctx context.Context, id int64) (*modelpayment.Payment, error) {
	if f.findOneFn == nil {
		panic("unexpected FindOne call")
	}
	return f.findOneFn(ctx, id)
}

func boolPtr(v bool) *bool {
	return &v
}

func newOrderTestDeps() Deps {
	return Deps{
		Config: &serverconfig.Config{},
	}
}

func requireCodeError(t *testing.T, err error, want uint32) {
	t.Helper()

	require.Error(t, err)

	var codeErr *xerr.CodeError
	require.ErrorAs(t, err, &codeErr)
	require.Equal(t, want, codeErr.GetErrCode())
}
