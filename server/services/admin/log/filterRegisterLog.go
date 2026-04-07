package log

import (
	"context"
	"github.com/perfect-panel/server/models/log"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type FilterRegisterLogInput struct {
	types.FilterRegisterLogRequest
}

type FilterRegisterLogOutput struct {
	Body *types.FilterRegisterLogResponse
}

func FilterRegisterLogHandler(svcCtx *svc.ServiceContext) func(context.Context, *FilterRegisterLogInput) (*FilterRegisterLogOutput, error) {
	return func(ctx context.Context, input *FilterRegisterLogInput) (*FilterRegisterLogOutput, error) {
		l := NewFilterRegisterLogLogic(ctx, svcCtx)
		resp, err := l.FilterRegisterLog(&input.FilterRegisterLogRequest)
		if err != nil {
			return nil, err
		}
		return &FilterRegisterLogOutput{Body: resp}, nil
	}
}

type FilterRegisterLogLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Filter register log
func NewFilterRegisterLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FilterRegisterLogLogic {
	return &FilterRegisterLogLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FilterRegisterLogLogic) FilterRegisterLog(req *types.FilterRegisterLogRequest) (resp *types.FilterRegisterLogResponse, err error) {
	data, total, err := l.svcCtx.LogModel.FilterSystemLog(l.ctx, &log.FilterParams{
		Page:     req.Page,
		Size:     req.Size,
		Type:     log.TypeRegister.Uint8(),
		ObjectID: req.UserId,
		Data:     req.Date,
		Search:   req.Search,
	})

	if err != nil {
		l.Errorf("[FilterRegisterLog] failed to filter system log: %v", err.Error())
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "failed to filter system log: %v", err.Error())
	}

	var list []types.RegisterLog
	for _, datum := range data {
		var item log.Register
		err = item.Unmarshal([]byte(datum.Content))
		if err != nil {
			l.Errorf("[FilterLoginLog] failed to unmarshal content: %v", err.Error())
			continue
		}
		list = append(list, types.RegisterLog{
			UserId:     datum.ObjectID,
			AuthMethod: item.AuthMethod,
			Identifier: item.Identifier,
			RegisterIP: item.RegisterIP,
			UserAgent:  item.UserAgent,
			Timestamp:  item.Timestamp,
		})
	}

	return &types.FilterRegisterLogResponse{
		List:  list,
		Total: total,
	}, nil
}
