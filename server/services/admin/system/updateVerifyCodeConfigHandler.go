// huma:migrated
package system

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type UpdateVerifyCodeConfigInput struct {
	Body types.VerifyCodeConfig
}

func UpdateVerifyCodeConfigHandler(svcCtx *svc.ServiceContext) func(context.Context, *UpdateVerifyCodeConfigInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateVerifyCodeConfigInput) (*struct{}, error) {
		l := NewUpdateVerifyCodeConfigLogic(ctx, svcCtx)
		if err := l.UpdateVerifyCodeConfig(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
