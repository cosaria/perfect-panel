package coupon

import (
	"context"

	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/pkg/errors"
)

type DeleteCouponInput struct {
	Body types.DeleteCouponRequest
}

func DeleteCouponHandler(deps Deps) func(context.Context, *DeleteCouponInput) (*struct{}, error) {
	return func(ctx context.Context, input *DeleteCouponInput) (*struct{}, error) {
		l := NewDeleteCouponLogic(ctx, deps)
		if err := l.DeleteCoupon(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type DeleteCouponLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Delete coupon
func NewDeleteCouponLogic(ctx context.Context, deps Deps) *DeleteCouponLogic {
	return &DeleteCouponLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *DeleteCouponLogic) DeleteCoupon(req *types.DeleteCouponRequest) error {
	// delete coupon by id
	err := l.deps.CouponModel.Delete(l.ctx, req.Id)
	if err != nil {
		l.Errorw("[DeleteCoupon] Database Error", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseDeletedError), "delete coupon error: %v", err.Error())
	}
	return nil
}
