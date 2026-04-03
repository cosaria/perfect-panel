// huma:migrated
package user

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/public/user"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type GetOAuthMethodsOutput struct {
	Body *types.GetOAuthMethodsResponse
}

func GetOAuthMethodsHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*GetOAuthMethodsOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetOAuthMethodsOutput, error) {
		l := user.NewGetOAuthMethodsLogic(ctx, svcCtx)
		resp, err := l.GetOAuthMethods()
		if err != nil {
			return nil, err
		}
		return &GetOAuthMethodsOutput{Body: resp}, nil
	}
}
