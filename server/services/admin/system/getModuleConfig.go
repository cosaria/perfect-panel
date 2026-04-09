package system

import (
	"context"
	"os"

	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/pkg/errors"
)

type GetModuleConfigOutput struct {
	Body *types.ModuleConfig
}

func GetModuleConfigHandler(deps Deps) func(context.Context, *struct{}) (*GetModuleConfigOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetModuleConfigOutput, error) {
		l := NewGetModuleConfigLogic(ctx, deps)
		resp, err := l.GetModuleConfig()
		if err != nil {
			return nil, err
		}
		return &GetModuleConfigOutput{Body: resp}, nil
	}
}

type GetModuleConfigLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Get Module Config
func NewGetModuleConfigLogic(ctx context.Context, deps Deps) *GetModuleConfigLogic {
	return &GetModuleConfigLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
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
