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

type GetRegisterConfigOutput struct {
	Body *types.RegisterConfig
}

func GetRegisterConfigHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*GetRegisterConfigOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetRegisterConfigOutput, error) {
		l := NewGetRegisterConfigLogic(ctx, svcCtx)
		resp, err := l.GetRegisterConfig()
		if err != nil {
			return nil, err
		}
		return &GetRegisterConfigOutput{Body: resp}, nil
	}
}

type GetRegisterConfigLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRegisterConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRegisterConfigLogic {
	return &GetRegisterConfigLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRegisterConfigLogic) GetRegisterConfig() (*types.RegisterConfig, error) {
	resp := &types.RegisterConfig{}

	// get register config from database
	configs, err := l.svcCtx.SystemModel.GetRegisterConfig(l.ctx)
	if err != nil {
		l.Errorw("[GetRegisterConfig] Database query error", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "get register config error: %v", err.Error())
	}

	// reflect to response
	tool.SystemConfigSliceReflectToStruct(configs, resp)
	return resp, nil
}
