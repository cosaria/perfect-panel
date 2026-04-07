// huma:migrated
package system

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
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
