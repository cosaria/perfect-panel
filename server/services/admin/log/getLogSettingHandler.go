// huma:migrated
package log

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type GetLogSettingOutput struct {
	Body *types.LogSetting
}

func GetLogSettingHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*GetLogSettingOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetLogSettingOutput, error) {
		l := NewGetLogSettingLogic(ctx, svcCtx)
		resp, err := l.GetLogSetting()
		if err != nil {
			return nil, err
		}
		return &GetLogSettingOutput{Body: resp}, nil
	}
}
