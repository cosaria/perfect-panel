package system

import (
	"context"
	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
	"os"
)

type GetModuleConfigOutput struct {
	Body *types.ModuleConfig
}

func GetModuleConfigHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*GetModuleConfigOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetModuleConfigOutput, error) {
		l := NewGetModuleConfigLogic(ctx, svcCtx)
		resp, err := l.GetModuleConfig()
		if err != nil {
			return nil, err
		}
		return &GetModuleConfigOutput{Body: resp}, nil
	}
}

type GetModuleConfigLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Get Module Config
func NewGetModuleConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetModuleConfigLogic {
	return &GetModuleConfigLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetModuleConfigLogic) GetModuleConfig() (resp *types.ModuleConfig, err error) {
	value, exists := os.LookupEnv("SECRET_KEY")
	if !exists {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.ERROR), " SECRET_KEY not set in environment variables")
	}

	return &types.ModuleConfig{
		Secret:         value,
		ServiceName:    config.ServiceName,
		ServiceVersion: config.Version,
	}, nil
}
