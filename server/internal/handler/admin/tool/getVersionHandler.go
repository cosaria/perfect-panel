// huma:migrated
package tool

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/tool"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type GetVersionOutput struct {
	Body *types.VersionResponse
}

func GetVersionHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*GetVersionOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetVersionOutput, error) {
		l := tool.NewGetVersionLogic(ctx, svcCtx)
		resp, err := l.GetVersion()
		if err != nil {
			return nil, err
		}
		return &GetVersionOutput{Body: resp}, nil
	}
}
