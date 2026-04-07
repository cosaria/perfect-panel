// huma:migrated
package user

import (
	"context"
	"github.com/perfect-panel/server/services/user/user"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type BindOAuthCallbackInput struct {
	Body types.BindOAuthCallbackRequest
}

func BindOAuthCallbackHandler(svcCtx *svc.ServiceContext) func(context.Context, *BindOAuthCallbackInput) (*struct{}, error) {
	return func(ctx context.Context, input *BindOAuthCallbackInput) (*struct{}, error) {
		l := user.NewBindOAuthCallbackLogic(ctx, svcCtx)
		if err := l.BindOAuthCallback(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
