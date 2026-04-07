package user

import (
	"context"

	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/pkg/errors"

	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type GetUserSubscribeTrafficLogsLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Get user subcribe traffic logs
func NewGetUserSubscribeTrafficLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserSubscribeTrafficLogsLogic {
	return &GetUserSubscribeTrafficLogsLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserSubscribeTrafficLogsLogic) GetUserSubscribeTrafficLogs(req *types.GetUserSubscribeTrafficLogsRequest) (resp *types.GetUserSubscribeTrafficLogsResponse, err error) {
	list, total, err := l.svcCtx.TrafficLogModel.QueryTrafficLogPageList(l.ctx, req.UserId, req.SubscribeId, req.Page, req.Size)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "GetUserSubscribeTrafficLogs failed: %v", err.Error())
	}
	userRespList := make([]types.TrafficLog, 0)
	tool.DeepCopy(&userRespList, list)
	return &types.GetUserSubscribeTrafficLogsResponse{
		Total: total,
		List:  userRespList,
	}, nil
}
