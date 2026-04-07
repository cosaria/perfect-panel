// huma:migrated
package user

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type UnbindOAuthInput struct {
	Body types.UnbindOAuthRequest
}

func UnbindOAuthHandler(svcCtx *svc.ServiceContext) func(context.Context, *UnbindOAuthInput) (*struct{}, error) {
	return func(ctx context.Context, input *UnbindOAuthInput) (*struct{}, error) {
		l := NewUnbindOAuthLogic(ctx, svcCtx)
		if err := l.UnbindOAuth(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
