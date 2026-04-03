// huma:migrated
package coupon

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/coupon"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type BatchDeleteCouponInput struct {
	Body types.BatchDeleteCouponRequest
}

func BatchDeleteCouponHandler(svcCtx *svc.ServiceContext) func(context.Context, *BatchDeleteCouponInput) (*struct{}, error) {
	return func(ctx context.Context, input *BatchDeleteCouponInput) (*struct{}, error) {
		l := coupon.NewBatchDeleteCouponLogic(ctx, svcCtx)
		if err := l.BatchDeleteCoupon(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
