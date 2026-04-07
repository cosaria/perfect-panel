// huma:migrated
package coupon

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type DeleteCouponInput struct {
	Body types.DeleteCouponRequest
}

func DeleteCouponHandler(svcCtx *svc.ServiceContext) func(context.Context, *DeleteCouponInput) (*struct{}, error) {
	return func(ctx context.Context, input *DeleteCouponInput) (*struct{}, error) {
		l := NewDeleteCouponLogic(ctx, svcCtx)
		if err := l.DeleteCoupon(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
