package coupon

import (
	"context"
	"math/rand"
	"time"

	"github.com/perfect-panel/server/internal/platform/crypto/snowflake"
	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/internal/platform/persistence/coupon"
	"github.com/perfect-panel/server/internal/platform/support/logger"
	"github.com/perfect-panel/server/internal/platform/support/random"
	"github.com/perfect-panel/server/internal/platform/support/tool"
	"github.com/perfect-panel/server/internal/platform/support/xerr"
	"github.com/pkg/errors"
)

type CreateCouponInput struct {
	Body types.CreateCouponRequest
}

func CreateCouponHandler(deps Deps) func(context.Context, *CreateCouponInput) (*struct{}, error) {
	return func(ctx context.Context, input *CreateCouponInput) (*struct{}, error) {
		l := NewCreateCouponLogic(ctx, deps)
		if err := l.CreateCoupon(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type CreateCouponLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Create coupon
func NewCreateCouponLogic(ctx context.Context, deps Deps) *CreateCouponLogic {
	return &CreateCouponLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *CreateCouponLogic) CreateCoupon(req *types.CreateCouponRequest) error {
	if req.Code == "" {
		rand.NewSource(time.Now().UnixNano())
		sid := snowflake.GetID()
		req.Code = random.KeyNew(4, 2) + "-" + random.StrToDashedString(random.EncodeBase36(sid))
	}
	couponInfo := &coupon.Coupon{}
	tool.DeepCopy(couponInfo, req)
	couponInfo.Subscribe = tool.Int64SliceToString(req.Subscribe)
	err := l.deps.CouponModel.Insert(l.ctx, couponInfo)
	if err != nil {
		l.Errorw("[CreateCoupon] Database Error", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseInsertError), "create coupon error: %v", err.Error())
	}
	return nil
}
