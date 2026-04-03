// huma:migrated
package user

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/public/user"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type UnbindDeviceInput struct {
	Body types.UnbindDeviceRequest
}

func UnbindDeviceHandler(svcCtx *svc.ServiceContext) func(context.Context, *UnbindDeviceInput) (*struct{}, error) {
	return func(ctx context.Context, input *UnbindDeviceInput) (*struct{}, error) {
		l := user.NewUnbindDeviceLogic(ctx, svcCtx)
		if err := l.UnbindDevice(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
