// huma:migrated
package application

import (
	"context"
	"github.com/perfect-panel/server/services/admin/application"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type UpdateSubscribeApplicationInput struct {
	Body types.UpdateSubscribeApplicationRequest
}

type UpdateSubscribeApplicationOutput struct {
	Body *types.SubscribeApplication
}

func UpdateSubscribeApplicationHandler(svcCtx *svc.ServiceContext) func(context.Context, *UpdateSubscribeApplicationInput) (*UpdateSubscribeApplicationOutput, error) {
	return func(ctx context.Context, input *UpdateSubscribeApplicationInput) (*UpdateSubscribeApplicationOutput, error) {
		l := application.NewUpdateSubscribeApplicationLogic(ctx, svcCtx)
		resp, err := l.UpdateSubscribeApplication(&input.Body)
		if err != nil {
			return nil, err
		}
		return &UpdateSubscribeApplicationOutput{Body: resp}, nil
	}
}
