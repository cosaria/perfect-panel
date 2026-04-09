package log

import (
	"context"

	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/internal/platform/persistence/log"
	"github.com/perfect-panel/server/internal/platform/support/logger"
	"github.com/perfect-panel/server/internal/platform/support/xerr"
	"github.com/pkg/errors"
)

type GetMessageLogListInput struct {
	types.GetMessageLogListRequest
}

type GetMessageLogListOutput struct {
	Body *types.GetMessageLogListResponse
}

func GetMessageLogListHandler(deps Deps) func(context.Context, *GetMessageLogListInput) (*GetMessageLogListOutput, error) {
	return func(ctx context.Context, input *GetMessageLogListInput) (*GetMessageLogListOutput, error) {
		l := NewGetMessageLogListLogic(ctx, deps)
		resp, err := l.GetMessageLogList(&input.GetMessageLogListRequest)
		if err != nil {
			return nil, err
		}
		return &GetMessageLogListOutput{Body: resp}, nil
	}
}

type GetMessageLogListLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// NewGetMessageLogListLogic Get message log list
func NewGetMessageLogListLogic(ctx context.Context, deps Deps) *GetMessageLogListLogic {
	return &GetMessageLogListLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *GetMessageLogListLogic) GetMessageLogList(req *types.GetMessageLogListRequest) (resp *types.GetMessageLogListResponse, err error) {

	data, total, err := l.deps.LogModel.FilterSystemLog(l.ctx, &log.FilterParams{
		Page:   req.Page,
		Size:   req.Size,
		Type:   req.Type,
		Search: req.Search,
	})

	if err != nil {
		l.Errorf("[GetMessageLogList] failed to filter system log: %v", err.Error())
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "failed to filter system log: %v", err.Error())
	}

	var list []types.MessageLog

	for _, datum := range data {
		var content log.Message
		err = content.Unmarshal([]byte(datum.Content))
		if err != nil {
			l.Errorf("[GetMessageLogList] failed to unmarshal content: %v", err.Error())
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

	return &types.GetMessageLogListResponse{
		Total: total,
		List:  list,
	}, nil
}
