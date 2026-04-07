// huma:migrated
package coupon

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type UpdateCouponInput struct {
	Body types.UpdateCouponRequest
}

func UpdateCouponHandler(svcCtx *svc.ServiceContext) func(context.Context, *UpdateCouponInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateCouponInput) (*struct{}, error) {
		l := NewUpdateCouponLogic(ctx, svcCtx)
		if err := l.UpdateCoupon(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
