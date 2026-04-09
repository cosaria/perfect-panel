package log

import (
	"context"
	"github.com/perfect-panel/server/models/log"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/pkg/errors"
)

type FilterLoginLogInput struct {
	types.FilterLoginLogRequest
}

type FilterLoginLogOutput struct {
	Body *types.FilterLoginLogResponse
}

func FilterLoginLogHandler(deps Deps) func(context.Context, *FilterLoginLogInput) (*FilterLoginLogOutput, error) {
	return func(ctx context.Context, input *FilterLoginLogInput) (*FilterLoginLogOutput, error) {
		l := NewFilterLoginLogLogic(ctx, deps)
		resp, err := l.FilterLoginLog(&input.FilterLoginLogRequest)
		if err != nil {
			return nil, err
		}
		return &FilterLoginLogOutput{Body: resp}, nil
	}
}

type FilterLoginLogLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// NewFilterLoginLogLogic Filter login log
func NewFilterLoginLogLogic(ctx context.Context, deps Deps) *FilterLoginLogLogic {
	return &FilterLoginLogLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *FilterLoginLogLogic) FilterLoginLog(req *types.FilterLoginLogRequest) (resp *types.FilterLoginLogResponse, err error) {
	data, total, err := l.deps.LogModel.FilterSystemLog(l.ctx, &log.FilterParams{
		Page:     req.Page,
		Size:     req.Size,
		Type:     log.TypeLogin.Uint8(),
		ObjectID: req.UserId,
		Data:     req.Date,
		Search:   req.Search,
	})

	if err != nil {
		l.Errorf("[FilterLoginLog] failed to filter system log: %v", err.Error())
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "failed to filter system log: %v", err.Error())
	}
	var list []types.LoginLog
	for _, datum := range data {
		var item log.Login
		err = item.Unmarshal([]byte(datum.Content))
		if err != nil {
			l.Errorf("[FilterLoginLog] failed to unmarshal content: %v", err.Error())
			continue
		}
		list = append(list, types.LoginLog{
			UserId:    datum.ObjectID,
			Method:    item.Method,
			LoginIP:   item.LoginIP,
			UserAgent: item.UserAgent,
			Success:   item.Success,
			Timestamp: datum.CreatedAt.UnixMilli(),
		})
	}

	return &types.FilterLoginLogResponse{
		Total: total,
		List:  list,
	}, nil
}
