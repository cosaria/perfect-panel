package system

import (
	"context"

	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/internal/platform/support/logger"
	"github.com/perfect-panel/server/internal/platform/support/tool"
	"github.com/perfect-panel/server/internal/platform/support/xerr"
	"github.com/pkg/errors"
)

type GetTosConfigOutput struct {
	Body *types.TosConfig
}

func GetTosConfigHandler(deps Deps) func(context.Context, *struct{}) (*GetTosConfigOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetTosConfigOutput, error) {
		l := NewGetTosConfigLogic(ctx, deps)
		resp, err := l.GetTosConfig()
		if err != nil {
			return nil, err
		}
		return &GetTosConfigOutput{Body: resp}, nil
	}
}

type GetTosConfigLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

func NewGetTosConfigLogic(ctx context.Context, deps Deps) *GetTosConfigLogic {
	return &GetTosConfigLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *GetTosConfigLogic) GetTosConfig() (resp *types.TosConfig, err error) {
	resp = &types.TosConfig{}
	// get tos config from db
	configs, err := l.deps.SystemModel.GetTosConfig(l.ctx)
	if err != nil {
		l.Errorw("[GetTosConfig] GetTosConfig error", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "GetTosConfig error: %v", err.Error())
	}
	// reflect to response
	tool.SystemConfigSliceReflectToStruct(configs, resp)
	return
}
