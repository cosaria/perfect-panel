// huma:migrated
package authMethod

import (
	"context"
	"github.com/perfect-panel/server/services/admin/authMethod"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type GetEmailPlatformOutput struct {
	Body *types.PlatformResponse
}

func GetEmailPlatformHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*GetEmailPlatformOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetEmailPlatformOutput, error) {
		l := authMethod.NewGetEmailPlatformLogic(ctx, svcCtx)
		resp, err := l.GetEmailPlatform()
		if err != nil {
			return nil, err
		}
		return &GetEmailPlatformOutput{Body: resp}, nil
	}
}
