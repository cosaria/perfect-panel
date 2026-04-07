package log

import (
	"context"
	"github.com/perfect-panel/server/models/log"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type FilterResetSubscribeLogInput struct {
	types.FilterResetSubscribeLogRequest
}

type FilterResetSubscribeLogOutput struct {
	Body *types.FilterResetSubscribeLogResponse
}

func FilterResetSubscribeLogHandler(deps Deps) func(context.Context, *FilterResetSubscribeLogInput) (*FilterResetSubscribeLogOutput, error) {
	return func(ctx context.Context, input *FilterResetSubscribeLogInput) (*FilterResetSubscribeLogOutput, error) {
		l := NewFilterResetSubscribeLogLogic(ctx, deps)
		resp, err := l.FilterResetSubscribeLog(&input.FilterResetSubscribeLogRequest)
		if err != nil {
			return nil, err
		}
		return &FilterResetSubscribeLogOutput{Body: resp}, nil
	}
}

type FilterResetSubscribeLogLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// NewFilterResetSubscribeLogLogic Filter reset subscribe log
func NewFilterResetSubscribeLogLogic(ctx context.Context, deps Deps) *FilterResetSubscribeLogLogic {
	return &FilterResetSubscribeLogLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *FilterResetSubscribeLogLogic) FilterResetSubscribeLog(req *types.FilterResetSubscribeLogRequest) (resp *types.FilterResetSubscribeLogResponse, err error) {
	data, total, err := l.deps.LogModel.FilterSystemLog(l.ctx, &log.FilterParams{
		Page:     req.Page,
		Size:     req.Size,
		Type:     log.TypeResetSubscribe.Uint8(),
		ObjectID: req.UserSubscribeId,
		Data:     req.Date,
		Search:   req.Search,
	})

	if err != nil {
		l.Errorf("[FilterResetSubscribeLog] failed to filter system log: %v", err.Error())
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "failed to filter system log: %v", err.Error())
	}

	var list []types.ResetSubscribeLog

	for _, item := range data {
		var content log.ResetSubscribe
		err = content.Unmarshal([]byte(item.Content))
		if err != nil {
			l.Errorf("[FilterResetSubscribeLog] failed to unmarshal content: %v", err.Error())
			continue
		}
		list = append(list, types.ResetSubscribeLog{
			Type:            content.Type,
			UserId:          content.UserId,
			UserSubscribeId: item.ObjectID,
			OrderNo:         content.OrderNo,
			Timestamp:       content.Timestamp,
		})
	}

	return &types.FilterResetSubscribeLogResponse{
		List:  list,
		Total: total,
	}, nil
}
