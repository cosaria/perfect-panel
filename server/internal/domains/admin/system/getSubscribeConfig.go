package system

import (
	"context"

	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/internal/platform/support/logger"
	"github.com/perfect-panel/server/internal/platform/support/tool"
	"github.com/perfect-panel/server/internal/platform/support/xerr"
	"github.com/pkg/errors"
)

type GetSubscribeConfigOutput struct {
	Body *types.SubscribeConfig
}

func GetSubscribeConfigHandler(deps Deps) func(context.Context, *struct{}) (*GetSubscribeConfigOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetSubscribeConfigOutput, error) {
		l := NewGetSubscribeConfigLogic(ctx, deps)
		resp, err := l.GetSubscribeConfig()
		if err != nil {
			return nil, err
		}
		return &GetSubscribeConfigOutput{Body: resp}, nil
	}
}

type GetSubscribeConfigLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

func NewGetSubscribeConfigLogic(ctx context.Context, deps Deps) *GetSubscribeConfigLogic {
	return &GetSubscribeConfigLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *GetSubscribeConfigLogic) GetSubscribeConfig() (resp *types.SubscribeConfig, err error) {
	resp = &types.SubscribeConfig{}
	// get subscribe config from db
	subscribeConfigs, err := l.deps.SystemModel.GetSubscribeConfig(l.ctx)
	if err != nil {
		l.Errorw("[GetSubscribeConfig] Database query error", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "get subscribe config failed: %v", err.Error())
	}

	// reflect to response
	tool.SystemConfigSliceReflectToStruct(subscribeConfigs, resp)
	return resp, nil
}
