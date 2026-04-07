package coupon

import (
	"context"
	"fmt"
	"github.com/perfect-panel/server/models/coupon"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
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

type UpdateCouponLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Update coupon
func NewUpdateCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCouponLogic {
	return &UpdateCouponLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateCouponLogic) UpdateCoupon(req *types.UpdateCouponRequest) error {
	fmt.Printf("req Subscribe: %v\n", req.Subscribe)
	couponInfo := &coupon.Coupon{}
	// update coupon
	tool.DeepCopy(couponInfo, req)
	couponInfo.Subscribe = tool.Int64SliceToString(req.Subscribe)
	err := l.svcCtx.CouponModel.Update(l.ctx, couponInfo)
	if err != nil {
		l.Errorw("[UpdateCoupon] Database Error", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseUpdateError), "update coupon error: %v", err.Error())
	}
	return nil
}
