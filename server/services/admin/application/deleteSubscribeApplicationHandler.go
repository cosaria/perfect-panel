// huma:migrated
package application

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type DeleteSubscribeApplicationInput struct {
	Body types.DeleteSubscribeApplicationRequest
}

func DeleteSubscribeApplicationHandler(svcCtx *svc.ServiceContext) func(context.Context, *DeleteSubscribeApplicationInput) (*struct{}, error) {
	return func(ctx context.Context, input *DeleteSubscribeApplicationInput) (*struct{}, error) {
		l := NewDeleteSubscribeApplicationLogic(ctx, svcCtx)
		if err := l.DeleteSubscribeApplication(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
