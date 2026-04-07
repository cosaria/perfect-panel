// huma:migrated
package system

import (
	"context"
	"github.com/perfect-panel/server/services/admin/system"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type GetVerifyCodeConfigOutput struct {
	Body *types.VerifyCodeConfig
}

func GetVerifyCodeConfigHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*GetVerifyCodeConfigOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetVerifyCodeConfigOutput, error) {
		l := system.NewGetVerifyCodeConfigLogic(ctx, svcCtx)
		resp, err := l.GetVerifyCodeConfig()
		if err != nil {
			return nil, err
		}
		return &GetVerifyCodeConfigOutput{Body: resp}, nil
	}
}
