package user

import (
	"context"
	"encoding/json"

	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/models/user"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/pkg/errors"
	"time"
)

type QueryUserSubscribeOutput struct {
	Body *types.QueryUserSubscribeListResponse
}

func QueryUserSubscribeHandler(deps Deps) func(context.Context, *struct{}) (*QueryUserSubscribeOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*QueryUserSubscribeOutput, error) {
		l := NewQueryUserSubscribeLogic(ctx, deps)
		resp, err := l.QueryUserSubscribe()
		if err != nil {
			return nil, err
		}
		return &QueryUserSubscribeOutput{Body: resp}, nil
	}
}

type QueryUserSubscribeLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Query User Subscribe
func NewQueryUserSubscribeLogic(ctx context.Context, deps Deps) *QueryUserSubscribeLogic {
	return &QueryUserSubscribeLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *QueryUserSubscribeLogic) QueryUserSubscribe() (resp *types.QueryUserSubscribeListResponse, err error) {
	u, ok := l.ctx.Value(config.CtxKeyUser).(*user.User)
	if !ok {
		logger.Error("current user is not found in context")
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.InvalidAccess), "Invalid Access")
	}
	data, err := l.deps.UserModel.QueryUserSubscribe(l.ctx, u.Id, 0, 1, 2, 3)
	if err != nil {
		l.Errorw("[QueryUserSubscribeLogic] Query User Subscribe Error:", logger.Field("err", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "Query User Subscribe Error")
	}

	resp = &types.QueryUserSubscribeListResponse{
		List:  make([]types.UserSubscribe, 0),
		Total: int64(len(data)),
	}

	for _, item := range data {
		var sub types.UserSubscribe
		tool.DeepCopy(&sub, item)

		// 解析Discount字段 避免在续订时只能续订一个月
		if item.Subscribe != nil && item.Subscribe.Discount != "" {
			var discounts []types.SubscribeDiscount
			if err := json.Unmarshal([]byte(item.Subscribe.Discount), &discounts); err == nil {
				sub.Subscribe.Discount = discounts
			}
		}

		short, _ := tool.FixedUniqueString(item.Token, 8, "")
		sub.Short = short
		sub.ResetTime = calculateNextResetTime(&sub)
		resp.List = append(resp.List, sub)
	}
	return
}

// 计算下次重置时间
func calculateNextResetTime(sub *types.UserSubscribe) int64 {
	resetTime := time.UnixMilli(sub.ExpireTime)
	now := time.Now()
	switch sub.Subscribe.ResetCycle {
	case 0:
		return 0
	case 1:
		return time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location()).UnixMilli()
	case 2:
		if resetTime.Day() > now.Day() {
			return time.Date(now.Year(), now.Month(), resetTime.Day(), 0, 0, 0, 0, now.Location()).UnixMilli()
		} else {
			return time.Date(now.Year(), now.Month()+1, resetTime.Day(), 0, 0, 0, 0, now.Location()).UnixMilli()
		}
	case 3:
		targetTime := time.Date(now.Year(), resetTime.Month(), resetTime.Day(), 0, 0, 0, 0, now.Location())
		if targetTime.Before(now) {
			targetTime = time.Date(now.Year()+1, resetTime.Month(), resetTime.Day(), 0, 0, 0, 0, now.Location())
		}
		return targetTime.UnixMilli()
	default:
		return 0
	}
}
