// huma:migrated
package tool

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type GetVersionOutput struct {
	Body *types.VersionResponse
}

func GetVersionHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*GetVersionOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetVersionOutput, error) {
		l := NewGetVersionLogic(ctx, svcCtx)
		resp, err := l.GetVersion()
		if err != nil {
			return nil, err
		}
		return &GetVersionOutput{Body: resp}, nil
	}
}
