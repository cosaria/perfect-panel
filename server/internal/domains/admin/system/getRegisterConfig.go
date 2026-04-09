package system

import (
	"context"

	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/internal/platform/support/logger"
	"github.com/perfect-panel/server/internal/platform/support/tool"
	"github.com/perfect-panel/server/internal/platform/support/xerr"
	"github.com/pkg/errors"
)

type GetRegisterConfigOutput struct {
	Body *types.RegisterConfig
}

func GetRegisterConfigHandler(deps Deps) func(context.Context, *struct{}) (*GetRegisterConfigOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetRegisterConfigOutput, error) {
		l := NewGetRegisterConfigLogic(ctx, deps)
		resp, err := l.GetRegisterConfig()
		if err != nil {
			return nil, err
		}
		return &GetRegisterConfigOutput{Body: resp}, nil
	}
}

type GetRegisterConfigLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

func NewGetRegisterConfigLogic(ctx context.Context, deps Deps) *GetRegisterConfigLogic {
	return &GetRegisterConfigLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *GetRegisterConfigLogic) GetRegisterConfig() (*types.RegisterConfig, error) {
	resp := &types.RegisterConfig{}

	// get register config from database
	configs, err := l.deps.SystemModel.GetRegisterConfig(l.ctx)
	if err != nil {
		l.Errorw("[GetRegisterConfig] Database query error", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "get register config error: %v", err.Error())
	}

	// reflect to response
	tool.SystemConfigSliceReflectToStruct(configs, resp)
	return resp, nil
}
