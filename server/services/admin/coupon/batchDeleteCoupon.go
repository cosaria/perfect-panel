package coupon

import (
	"context"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/pkg/errors"
)

type BatchDeleteCouponInput struct {
	Body types.BatchDeleteCouponRequest
}

func BatchDeleteCouponHandler(deps Deps) func(context.Context, *BatchDeleteCouponInput) (*struct{}, error) {
	return func(ctx context.Context, input *BatchDeleteCouponInput) (*struct{}, error) {
		l := NewBatchDeleteCouponLogic(ctx, deps)
		if err := l.BatchDeleteCoupon(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type BatchDeleteCouponLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Batch delete coupon
func NewBatchDeleteCouponLogic(ctx context.Context, deps Deps) *BatchDeleteCouponLogic {
	return &BatchDeleteCouponLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *BatchDeleteCouponLogic) BatchDeleteCoupon(req *types.BatchDeleteCouponRequest) error {
	// batch delete coupon by ids
	err := l.deps.CouponModel.BatchDelete(l.ctx, req.Ids)
	if err != nil {
		l.Errorw("[BatchDeleteCoupon] Database Error", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseDeletedError), "batch delete coupon error: %v", err.Error())
	}
	return nil
}
