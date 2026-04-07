// huma:migrated
package authMethod

import (
	"context"
	"github.com/perfect-panel/server/services/admin/authMethod"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type UpdateAuthMethodConfigInput struct {
	Body types.UpdateAuthMethodConfigRequest
}

type UpdateAuthMethodConfigOutput struct {
	Body *types.AuthMethodConfig
}

func UpdateAuthMethodConfigHandler(svcCtx *svc.ServiceContext) func(context.Context, *UpdateAuthMethodConfigInput) (*UpdateAuthMethodConfigOutput, error) {
	return func(ctx context.Context, input *UpdateAuthMethodConfigInput) (*UpdateAuthMethodConfigOutput, error) {
		l := authMethod.NewUpdateAuthMethodConfigLogic(ctx, svcCtx)
		resp, err := l.UpdateAuthMethodConfig(&input.Body)
		if err != nil {
			return nil, err
		}
		return &UpdateAuthMethodConfigOutput{Body: resp}, nil
	}
}
