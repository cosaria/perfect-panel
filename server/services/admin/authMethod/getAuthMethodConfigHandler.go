// huma:migrated
package authMethod

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type GetAuthMethodConfigInput struct {
	types.GetAuthMethodConfigRequest
}

type GetAuthMethodConfigOutput struct {
	Body *types.AuthMethodConfig
}

func GetAuthMethodConfigHandler(svcCtx *svc.ServiceContext) func(context.Context, *GetAuthMethodConfigInput) (*GetAuthMethodConfigOutput, error) {
	return func(ctx context.Context, input *GetAuthMethodConfigInput) (*GetAuthMethodConfigOutput, error) {
		l := NewGetAuthMethodConfigLogic(ctx, svcCtx)
		resp, err := l.GetAuthMethodConfig(&input.GetAuthMethodConfigRequest)
		if err != nil {
			return nil, err
		}
		return &GetAuthMethodConfigOutput{Body: resp}, nil
	}
}
