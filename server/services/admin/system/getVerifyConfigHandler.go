// huma:migrated
package system

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type GetVerifyConfigOutput struct {
	Body *types.VerifyConfig
}

func GetVerifyConfigHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*GetVerifyConfigOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetVerifyConfigOutput, error) {
		l := NewGetVerifyConfigLogic(ctx, svcCtx)
		resp, err := l.GetVerifyConfig()
		if err != nil {
			return nil, err
		}
		return &GetVerifyConfigOutput{Body: resp}, nil
	}
}
