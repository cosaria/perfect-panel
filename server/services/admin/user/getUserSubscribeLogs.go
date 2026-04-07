package user

import (
	"context"
	"github.com/perfect-panel/server/models/log"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type GetUserSubscribeLogsInput struct {
	types.GetUserSubscribeLogsRequest
}

type GetUserSubscribeLogsOutput struct {
	Body *types.GetUserSubscribeLogsResponse
}

func GetUserSubscribeLogsHandler(svcCtx *svc.ServiceContext) func(context.Context, *GetUserSubscribeLogsInput) (*GetUserSubscribeLogsOutput, error) {
	return func(ctx context.Context, input *GetUserSubscribeLogsInput) (*GetUserSubscribeLogsOutput, error) {
		l := NewGetUserSubscribeLogsLogic(ctx, svcCtx)
		resp, err := l.GetUserSubscribeLogs(&input.GetUserSubscribeLogsRequest)
		if err != nil {
			return nil, err
		}
		return &GetUserSubscribeLogsOutput{Body: resp}, nil
	}
}

type GetUserSubscribeLogsLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Get user subcribe logs
func NewGetUserSubscribeLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserSubscribeLogsLogic {
	return &GetUserSubscribeLogsLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserSubscribeLogsLogic) GetUserSubscribeLogs(req *types.GetUserSubscribeLogsRequest) (resp *types.GetUserSubscribeLogsResponse, err error) {
	data, total, err := l.svcCtx.LogModel.FilterSystemLog(l.ctx, &log.FilterParams{})

	if err != nil {
		l.Errorw("[GetUserSubscribeLogs] Get User Subscribe Logs Error:", logger.Field("err", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "Get User Subscribe Logs Error")
	}
	var list []types.UserSubscribeLog
	tool.DeepCopy(&list, data)

	return &types.GetUserSubscribeLogsResponse{
		List:  list,
		Total: total,
	}, err
}
