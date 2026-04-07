package user

import (
	"context"

	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/models/log"
	"github.com/perfect-panel/server/models/user"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type QueryUserCommissionLogInput struct {
	types.QueryUserCommissionLogListRequest
}

type QueryUserCommissionLogOutput struct {
	Body *types.QueryUserCommissionLogListResponse
}

func QueryUserCommissionLogHandler(deps Deps) func(context.Context, *QueryUserCommissionLogInput) (*QueryUserCommissionLogOutput, error) {
	return func(ctx context.Context, input *QueryUserCommissionLogInput) (*QueryUserCommissionLogOutput, error) {
		l := NewQueryUserCommissionLogLogic(ctx, deps)
		resp, err := l.QueryUserCommissionLog(&input.QueryUserCommissionLogListRequest)
		if err != nil {
			return nil, err
		}
		return &QueryUserCommissionLogOutput{Body: resp}, nil
	}
}

type QueryUserCommissionLogLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Query User Commission Log
func NewQueryUserCommissionLogLogic(ctx context.Context, deps Deps) *QueryUserCommissionLogLogic {
	return &QueryUserCommissionLogLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *QueryUserCommissionLogLogic) QueryUserCommissionLog(req *types.QueryUserCommissionLogListRequest) (resp *types.QueryUserCommissionLogListResponse, err error) {
	u, ok := l.ctx.Value(config.CtxKeyUser).(*user.User)
	if !ok {
		logger.Error("current user is not found in context")
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.InvalidAccess), "Invalid Access")
	}
	data, total, err := l.deps.LogModel.FilterSystemLog(l.ctx, &log.FilterParams{
		Page:     req.Page,
		Size:     req.Size,
		Type:     log.TypeCommission.Uint8(),
		ObjectID: u.Id,
	})
	if err != nil {
		l.Errorw("Query User Commission Log failed", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "Query User Commission Log failed: %v", err)
	}
	var list []types.CommissionLog

	for _, datum := range data {
		var content log.Commission
		if err = content.Unmarshal([]byte(datum.Content)); err != nil {
			l.Errorf("unmarshal commission log content failed: %v", err.Error())
			continue
		}
		list = append(list, types.CommissionLog{
			UserId:    datum.ObjectID,
			Type:      content.Type,
			Amount:    content.Amount,
			OrderNo:   content.OrderNo,
			Timestamp: content.Timestamp,
		})
	}

	return &types.QueryUserCommissionLogListResponse{
		List:  list,
		Total: total,
	}, nil
}
