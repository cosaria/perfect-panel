package user

import (
	"context"

	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/models/log"
	"github.com/perfect-panel/server/models/user"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/pkg/errors"
)

type QueryUserBalanceLogOutput struct {
	Body *types.QueryUserBalanceLogListResponse
}

func QueryUserBalanceLogHandler(deps Deps) func(context.Context, *struct{}) (*QueryUserBalanceLogOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*QueryUserBalanceLogOutput, error) {
		l := NewQueryUserBalanceLogLogic(ctx, deps)
		resp, err := l.QueryUserBalanceLog()
		if err != nil {
			return nil, err
		}
		return &QueryUserBalanceLogOutput{Body: resp}, nil
	}
}

type QueryUserBalanceLogLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// NewQueryUserBalanceLogLogic Query User Balance Log
func NewQueryUserBalanceLogLogic(ctx context.Context, deps Deps) *QueryUserBalanceLogLogic {
	return &QueryUserBalanceLogLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *QueryUserBalanceLogLogic) QueryUserBalanceLog() (resp *types.QueryUserBalanceLogListResponse, err error) {
	u, ok := l.ctx.Value(config.CtxKeyUser).(*user.User)
	if !ok {
		logger.Error("current user is not found in context")
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.InvalidAccess), "Invalid Access")
	}

	data, total, err := l.deps.LogModel.FilterSystemLog(l.ctx, &log.FilterParams{
		Page:     1,
		Size:     99999,
		Type:     log.TypeBalance.Uint8(),
		ObjectID: u.Id,
	})
	if err != nil {
		l.Errorw("[QueryUserBalanceLog] Query User Balance Log Error:", logger.Field("err", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "Query User Balance Log Error")
	}

	list := make([]types.BalanceLog, 0)
	for _, datum := range data {
		var content log.Balance
		if err = content.Unmarshal([]byte(datum.Content)); err != nil {
			l.Errorf("[QueryUserBalanceLog] unmarshal balance log content failed: %v", err.Error())
			continue
		}
		list = append(list, types.BalanceLog{
			UserId:    datum.ObjectID,
			Amount:    content.Amount,
			Type:      content.Type,
			OrderNo:   content.OrderNo,
			Balance:   content.Balance,
			Timestamp: content.Timestamp,
		})
	}

	return &types.QueryUserBalanceLogListResponse{
		Total: total,
		List:  list,
	}, nil
}
