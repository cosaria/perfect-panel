// huma:migrated
package user

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/user"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type GetUserAuthMethodInput struct {
	types.GetUserAuthMethodRequest
}

type GetUserAuthMethodOutput struct {
	Body *types.GetUserAuthMethodResponse
}

func GetUserAuthMethodHandler(svcCtx *svc.ServiceContext) func(context.Context, *GetUserAuthMethodInput) (*GetUserAuthMethodOutput, error) {
	return func(ctx context.Context, input *GetUserAuthMethodInput) (*GetUserAuthMethodOutput, error) {
		l := user.NewGetUserAuthMethodLogic(ctx, svcCtx)
		resp, err := l.GetUserAuthMethod(&input.GetUserAuthMethodRequest)
		if err != nil {
			return nil, err
		}
		return &GetUserAuthMethodOutput{Body: resp}, nil
	}
}
