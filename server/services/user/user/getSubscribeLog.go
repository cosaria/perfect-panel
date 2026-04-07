package user

import (
	"context"
	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/models/log"
	"github.com/perfect-panel/server/models/user"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type GetSubscribeLogInput struct {
	types.GetSubscribeLogRequest
}

type GetSubscribeLogOutput struct {
	Body *types.GetSubscribeLogResponse
}

func GetSubscribeLogHandler(svcCtx *svc.ServiceContext) func(context.Context, *GetSubscribeLogInput) (*GetSubscribeLogOutput, error) {
	return func(ctx context.Context, input *GetSubscribeLogInput) (*GetSubscribeLogOutput, error) {
		l := NewGetSubscribeLogLogic(ctx, svcCtx)
		resp, err := l.GetSubscribeLog(&input.GetSubscribeLogRequest)
		if err != nil {
			return nil, err
		}
		return &GetSubscribeLogOutput{Body: resp}, nil
	}
}

type GetSubscribeLogLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetSubscribeLogLogic Get Subscribe Log
func NewGetSubscribeLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSubscribeLogLogic {
	return &GetSubscribeLogLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSubscribeLogLogic) GetSubscribeLog(req *types.GetSubscribeLogRequest) (resp *types.GetSubscribeLogResponse, err error) {
	u, ok := l.ctx.Value(config.CtxKeyUser).(*user.User)
	if !ok {
		logger.Error("current user is not found in context")
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.InvalidAccess), "Invalid Access")
	}
	data, total, err := l.svcCtx.LogModel.FilterSystemLog(l.ctx, &log.FilterParams{
		Page:     req.Page,
		Size:     req.Size,
		Type:     log.TypeSubscribe.Uint8(),
		ObjectID: u.Id, // filter by current user id
	})
	if err != nil {
		l.Errorw("[GetUserSubscribeLogs] Get User Subscribe Logs Error:", logger.Field("err", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "Get User Subscribe Logs Error")
	}
	var list []types.UserSubscribeLog

	for _, item := range data {
		var content log.Subscribe
		if err = content.Unmarshal([]byte(item.Content)); err != nil {
			l.Errorf("[GetUserSubscribeLogs] unmarshal subscribe log content failed: %v", err.Error())
			continue
		}
		list = append(list, types.UserSubscribeLog{
			Id:              item.Id,
			UserId:          item.ObjectID,
			UserSubscribeId: content.UserSubscribeId,
			Token:           content.Token,
			IP:              content.ClientIP,
			UserAgent:       content.UserAgent,
			Timestamp:       item.CreatedAt.UnixMilli(),
		})
	}

	return &types.GetSubscribeLogResponse{
		List:  list,
		Total: total,
	}, err
}
