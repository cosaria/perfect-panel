package log

import (
	"context"
	"strconv"

	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/models/log"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/pkg/errors"
)

type FilterSubscribeLogInput struct {
	types.FilterSubscribeLogRequest
}

type FilterSubscribeLogOutput struct {
	Body *types.FilterSubscribeLogResponse
}

func FilterSubscribeLogHandler(deps Deps) func(context.Context, *FilterSubscribeLogInput) (*FilterSubscribeLogOutput, error) {
	return func(ctx context.Context, input *FilterSubscribeLogInput) (*FilterSubscribeLogOutput, error) {
		l := NewFilterSubscribeLogLogic(ctx, deps)
		resp, err := l.FilterSubscribeLog(&input.FilterSubscribeLogRequest)
		if err != nil {
			return nil, err
		}
		return &FilterSubscribeLogOutput{Body: resp}, nil
	}
}

type FilterSubscribeLogLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// NewFilterSubscribeLogLogic Filter subscribe log
func NewFilterSubscribeLogLogic(ctx context.Context, deps Deps) *FilterSubscribeLogLogic {
	return &FilterSubscribeLogLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *FilterSubscribeLogLogic) FilterSubscribeLog(req *types.FilterSubscribeLogRequest) (resp *types.FilterSubscribeLogResponse, err error) {
	params := &log.FilterParams{
		Page:     req.Page,
		Size:     req.Size,
		Type:     log.TypeSubscribe.Uint8(),
		Data:     req.Date,
		ObjectID: req.UserId,
	}

	if req.UserSubscribeId != 0 {
		params.Search = `"user_subscribe_id":` + strconv.FormatInt(req.UserSubscribeId, 10)
	}

	data, total, err := l.deps.LogModel.FilterSystemLog(l.ctx, params)
	if err != nil {
		l.Errorf("[FilterSubscribeLog] failed to filter system log: %v", err.Error())
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "failed to filter system log")
	}

	var list []types.SubscribeLog
	for _, datum := range data {
		var content log.Subscribe
		err = content.Unmarshal([]byte(datum.Content))
		if err != nil {
			l.Errorf("[FilterSubscribeLog] failed to unmarshal content: %v", err.Error())
			continue
		}
		list = append(list, types.SubscribeLog{
			UserId:          datum.ObjectID,
			Token:           content.Token,
			UserAgent:       content.UserAgent,
			ClientIP:        content.ClientIP,
			UserSubscribeId: content.UserSubscribeId,
			Timestamp:       datum.CreatedAt.UnixMilli(),
		})
	}

	return &types.FilterSubscribeLogResponse{
		Total: total,
		List:  list,
	}, nil
}
