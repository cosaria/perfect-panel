package log

import (
	"context"

	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/util/tool"
)

type GetLogSettingOutput struct {
	Body *types.LogSetting
}

func GetLogSettingHandler(deps Deps) func(context.Context, *struct{}) (*GetLogSettingOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetLogSettingOutput, error) {
		l := NewGetLogSettingLogic(ctx, deps)
		resp, err := l.GetLogSetting()
		if err != nil {
			return nil, err
		}
		return &GetLogSettingOutput{Body: resp}, nil
	}
}

type GetLogSettingLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Get log setting
func NewGetLogSettingLogic(ctx context.Context, deps Deps) *GetLogSettingLogic {
	return &GetLogSettingLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *GetLogSettingLogic) GetLogSetting() (resp *types.LogSetting, err error) {
	configs, err := l.deps.SystemModel.GetLogConfig(l.ctx)
	if err != nil {
		l.Errorw("[GetLogSetting] Database query error", logger.Field("error", err.Error()))
		return nil, err
	}
	resp = &types.LogSetting{}
	// reflect to response
	tool.SystemConfigSliceReflectToStruct(configs, resp)
	return
}
