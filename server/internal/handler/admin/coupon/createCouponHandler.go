// huma:migrated
package coupon

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/coupon"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type CreateCouponInput struct {
	Body types.CreateCouponRequest
}

func CreateCouponHandler(svcCtx *svc.ServiceContext) func(context.Context, *CreateCouponInput) (*struct{}, error) {
	return func(ctx context.Context, input *CreateCouponInput) (*struct{}, error) {
		l := coupon.NewCreateCouponLogic(ctx, svcCtx)
		if err := l.CreateCoupon(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
