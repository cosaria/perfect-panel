// huma:migrated
package auth

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/auth"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type DeviceLoginInput struct {
	Body types.DeviceLoginRequest
}

type DeviceLoginOutput struct {
	Body *types.LoginResponse
}

func DeviceLoginHandler(svcCtx *svc.ServiceContext) func(context.Context, *DeviceLoginInput) (*DeviceLoginOutput, error) {
	return func(ctx context.Context, input *DeviceLoginInput) (*DeviceLoginOutput, error) {
		l := auth.NewDeviceLoginLogic(ctx, svcCtx)
		resp, err := l.DeviceLogin(&input.Body)
		if err != nil {
			return nil, err
		}
		return &DeviceLoginOutput{Body: resp}, nil
	}
}
