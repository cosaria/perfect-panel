// huma:migrated
package authMethod

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/authMethod"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type GetSmsPlatformOutput struct {
	Body *types.PlatformResponse
}

func GetSmsPlatformHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*GetSmsPlatformOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetSmsPlatformOutput, error) {
		l := authMethod.NewGetSmsPlatformLogic(ctx, svcCtx)
		resp, err := l.GetSmsPlatform()
		if err != nil {
			return nil, err
		}
		return &GetSmsPlatformOutput{Body: resp}, nil
	}
}
