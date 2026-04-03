// huma:migrated
package system

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/system"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type GetVerifyConfigOutput struct {
	Body *types.VerifyConfig
}

func GetVerifyConfigHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*GetVerifyConfigOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetVerifyConfigOutput, error) {
		l := system.NewGetVerifyConfigLogic(ctx, svcCtx)
		resp, err := l.GetVerifyConfig()
		if err != nil {
			return nil, err
		}
		return &GetVerifyConfigOutput{Body: resp}, nil
	}
}
