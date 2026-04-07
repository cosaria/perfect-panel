package system

import (
	"context"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type GetSubscribeConfigOutput struct {
	Body *types.SubscribeConfig
}

func GetSubscribeConfigHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*GetSubscribeConfigOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetSubscribeConfigOutput, error) {
		l := NewGetSubscribeConfigLogic(ctx, svcCtx)
		resp, err := l.GetSubscribeConfig()
		if err != nil {
			return nil, err
		}
		return &GetSubscribeConfigOutput{Body: resp}, nil
	}
}

type GetSubscribeConfigLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSubscribeConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSubscribeConfigLogic {
	return &GetSubscribeConfigLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSubscribeConfigLogic) GetSubscribeConfig() (resp *types.SubscribeConfig, err error) {
	resp = &types.SubscribeConfig{}
	// get subscribe config from db
	subscribeConfigs, err := l.svcCtx.SystemModel.GetSubscribeConfig(l.ctx)
	if err != nil {
		l.Errorw("[GetSubscribeConfig] Database query error", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "get subscribe config failed: %v", err.Error())
	}

	// reflect to response
	tool.SystemConfigSliceReflectToStruct(subscribeConfigs, resp)
	return resp, nil
}
