// huma:migrated
package user

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type KickOfflineByUserDeviceInput struct {
	Body types.KickOfflineRequest
}

func KickOfflineByUserDeviceHandler(svcCtx *svc.ServiceContext) func(context.Context, *KickOfflineByUserDeviceInput) (*struct{}, error) {
	return func(ctx context.Context, input *KickOfflineByUserDeviceInput) (*struct{}, error) {
		l := NewKickOfflineByUserDeviceLogic(ctx, svcCtx)
		if err := l.KickOfflineByUserDevice(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
