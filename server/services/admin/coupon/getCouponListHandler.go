// huma:migrated
package coupon

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type GetCouponListInput struct {
	types.GetCouponListRequest
}

type GetCouponListOutput struct {
	Body *types.GetCouponListResponse
}

func GetCouponListHandler(svcCtx *svc.ServiceContext) func(context.Context, *GetCouponListInput) (*GetCouponListOutput, error) {
	return func(ctx context.Context, input *GetCouponListInput) (*GetCouponListOutput, error) {
		l := NewGetCouponListLogic(ctx, svcCtx)
		resp, err := l.GetCouponList(&input.GetCouponListRequest)
		if err != nil {
			return nil, err
		}
		return &GetCouponListOutput{Body: resp}, nil
	}
}
