package log

import (
	"context"
	"github.com/perfect-panel/server/models/log"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type FilterEmailLogInput struct {
	types.FilterLogParams
}

type FilterEmailLogOutput struct {
	Body *types.FilterEmailLogResponse
}

func FilterEmailLogHandler(deps Deps) func(context.Context, *FilterEmailLogInput) (*FilterEmailLogOutput, error) {
	return func(ctx context.Context, input *FilterEmailLogInput) (*FilterEmailLogOutput, error) {
		l := NewFilterEmailLogLogic(ctx, deps)
		resp, err := l.FilterEmailLog(&input.FilterLogParams)
		if err != nil {
			return nil, err
		}
		return &FilterEmailLogOutput{Body: resp}, nil
	}
}

type FilterEmailLogLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// NewFilterEmailLogLogic Filter email log
func NewFilterEmailLogLogic(ctx context.Context, deps Deps) *FilterEmailLogLogic {
	return &FilterEmailLogLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *FilterEmailLogLogic) FilterEmailLog(req *types.FilterLogParams) (resp *types.FilterEmailLogResponse, err error) {
	data, total, err := l.deps.LogModel.FilterSystemLog(l.ctx, &log.FilterParams{
		Page:   req.Page,
		Size:   req.Size,
		Type:   log.TypeEmailMessage.Uint8(),
		Data:   req.Date,
		Search: req.Search,
	})

	if err != nil {
		l.Errorf("[FilterEmailLog] failed to filter system log: %v", err.Error())
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "failed to filter system log: %v", err.Error())
	}

	var list []types.MessageLog

	for _, datum := range data {
		var content log.Message
		err = content.Unmarshal([]byte(datum.Content))
		if err != nil {
			l.Errorf("[FilterEmailLog] failed to unmarshal content: %v", err.Error())
			continue
		}
		list = append(list, types.MessageLog{
			Id:        datum.Id,
			Type:      datum.Type,
			Platform:  content.Platform,
			To:        content.To,
			Subject:   content.Subject,
			Content:   content.Content,
			Status:    content.Status,
			CreatedAt: datum.CreatedAt.UnixMilli(),
		})
	}

	return &types.FilterEmailLogResponse{
		Total: total,
		List:  list,
	}, nil
}
