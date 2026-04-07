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

type GetTosConfigOutput struct {
	Body *types.TosConfig
}

func GetTosConfigHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*GetTosConfigOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetTosConfigOutput, error) {
		l := NewGetTosConfigLogic(ctx, svcCtx)
		resp, err := l.GetTosConfig()
		if err != nil {
			return nil, err
		}
		return &GetTosConfigOutput{Body: resp}, nil
	}
}

type GetTosConfigLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTosConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTosConfigLogic {
	return &GetTosConfigLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTosConfigLogic) GetTosConfig() (resp *types.TosConfig, err error) {
	resp = &types.TosConfig{}
	// get tos config from db
	configs, err := l.svcCtx.SystemModel.GetTosConfig(l.ctx)
	if err != nil {
		l.Errorw("[GetTosConfig] GetTosConfig error", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "GetTosConfig error: %v", err.Error())
	}
	// reflect to response
	tool.SystemConfigSliceReflectToStruct(configs, resp)
	return
}
