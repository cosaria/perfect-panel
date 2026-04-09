package system

import (
	"context"

	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/pkg/errors"
)

type GetCurrencyConfigOutput struct {
	Body *types.CurrencyConfig
}

func GetCurrencyConfigHandler(deps Deps) func(context.Context, *struct{}) (*GetCurrencyConfigOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetCurrencyConfigOutput, error) {
		l := NewGetCurrencyConfigLogic(ctx, deps)
		resp, err := l.GetCurrencyConfig()
		if err != nil {
			return nil, err
		}
		return &GetCurrencyConfigOutput{Body: resp}, nil
	}
}

type GetCurrencyConfigLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Get Currency Config
func NewGetCurrencyConfigLogic(ctx context.Context, deps Deps) *GetCurrencyConfigLogic {
	return &GetCurrencyConfigLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *GetCurrencyConfigLogic) GetCurrencyConfig() (resp *types.CurrencyConfig, err error) {
	configs, err := l.deps.SystemModel.GetCurrencyConfig(l.ctx)
	if err != nil {
		l.Errorw("[GetCurrencyConfigLogic] GetCurrencyConfig error: ", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "GetCurrencyConfig error: %v", err.Error())
	}
	resp = &types.CurrencyConfig{}
	tool.SystemConfigSliceReflectToStruct(configs, resp)
	return
}
