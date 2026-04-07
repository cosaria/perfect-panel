// huma:migrated
package auth

import (
	"context"
	"github.com/perfect-panel/server/services/auth"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type DeviceLoginInput struct {
	Body types.DeviceLoginRequest
	IP   string `header:"X-Original-Forwarded-For" required:"false" doc:"Client IP from proxy"`
}

type DeviceLoginOutput struct {
	Body *types.LoginResponse
}

func DeviceLoginHandler(svcCtx *svc.ServiceContext) func(context.Context, *DeviceLoginInput) (*DeviceLoginOutput, error) {
	return func(ctx context.Context, input *DeviceLoginInput) (*DeviceLoginOutput, error) {
		input.Body.IP = input.IP
		l := auth.NewDeviceLoginLogic(ctx, svcCtx)
		resp, err := l.DeviceLogin(&input.Body)
		if err != nil {
			return nil, err
		}
		return &DeviceLoginOutput{Body: resp}, nil
	}
}
