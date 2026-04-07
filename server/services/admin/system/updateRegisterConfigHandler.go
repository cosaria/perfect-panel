// huma:migrated
package system

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type UpdateRegisterConfigInput struct {
	Body types.RegisterConfig
}

func UpdateRegisterConfigHandler(svcCtx *svc.ServiceContext) func(context.Context, *UpdateRegisterConfigInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateRegisterConfigInput) (*struct{}, error) {
		l := NewUpdateRegisterConfigLogic(ctx, svcCtx)
		if err := l.UpdateRegisterConfig(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
