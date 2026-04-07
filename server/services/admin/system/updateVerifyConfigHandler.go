// huma:migrated
package system

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type UpdateVerifyConfigInput struct {
	Body types.VerifyConfig
}

func UpdateVerifyConfigHandler(svcCtx *svc.ServiceContext) func(context.Context, *UpdateVerifyConfigInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateVerifyConfigInput) (*struct{}, error) {
		l := NewUpdateVerifyConfigLogic(ctx, svcCtx)
		if err := l.UpdateVerifyConfig(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
