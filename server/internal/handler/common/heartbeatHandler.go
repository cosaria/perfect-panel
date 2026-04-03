// huma:migrated
package common

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/common"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type HeartbeatOutput struct {
	Body *types.HeartbeatResponse
}

func HeartbeatHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*HeartbeatOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*HeartbeatOutput, error) {
		l := common.NewHeartbeatLogic(ctx, svcCtx)
		resp, err := l.Heartbeat()
		if err != nil {
			return nil, err
		}
		return &HeartbeatOutput{Body: resp}, nil
	}
}
