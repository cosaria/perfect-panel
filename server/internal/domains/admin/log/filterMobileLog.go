package log

import (
	"context"

	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/internal/platform/persistence/log"
	"github.com/perfect-panel/server/internal/platform/support/logger"
	"github.com/perfect-panel/server/internal/platform/support/xerr"
	"github.com/pkg/errors"
)

type FilterMobileLogInput struct {
	types.FilterLogParams
}

type FilterMobileLogOutput struct {
	Body *types.FilterMobileLogResponse
}

func FilterMobileLogHandler(deps Deps) func(context.Context, *FilterMobileLogInput) (*FilterMobileLogOutput, error) {
	return func(ctx context.Context, input *FilterMobileLogInput) (*FilterMobileLogOutput, error) {
		l := NewFilterMobileLogLogic(ctx, deps)
		resp, err := l.FilterMobileLog(&input.FilterLogParams)
		if err != nil {
			return nil, err
		}
		return &FilterMobileLogOutput{Body: resp}, nil
	}
}

type FilterMobileLogLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Filter mobile log
func NewFilterMobileLogLogic(ctx context.Context, deps Deps) *FilterMobileLogLogic {
	return &FilterMobileLogLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *FilterMobileLogLogic) FilterMobileLog(req *types.FilterLogParams) (resp *types.FilterMobileLogResponse, err error) {
	data, total, err := l.deps.LogModel.FilterSystemLog(l.ctx, &log.FilterParams{
		Page:   req.Page,
		Size:   req.Size,
		Type:   log.TypeMobileMessage.Uint8(),
		Data:   req.Date,
		Search: req.Search,
	})

	if err != nil {
		l.Errorf("[FilterMobileLog] failed to filter system log: %v", err.Error())
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "failed to filter system log: %v", err.Error())
	}

	var list []types.MessageLog

	for _, datum := range data {
		var content log.Message
		err = content.Unmarshal([]byte(datum.Content))
		if err != nil {
			l.Errorf("[FilterMobileLog] failed to unmarshal content: %v", err.Error())
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

	return &types.FilterMobileLogResponse{
		Total: total,
		List:  list,
	}, nil
}
