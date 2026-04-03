// huma:migrated
package application

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/application"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type CreateSubscribeApplicationInput struct {
	Body types.CreateSubscribeApplicationRequest
}

type CreateSubscribeApplicationOutput struct {
	Body *types.SubscribeApplication
}

func CreateSubscribeApplicationHandler(svcCtx *svc.ServiceContext) func(context.Context, *CreateSubscribeApplicationInput) (*CreateSubscribeApplicationOutput, error) {
	return func(ctx context.Context, input *CreateSubscribeApplicationInput) (*CreateSubscribeApplicationOutput, error) {
		l := application.NewCreateSubscribeApplicationLogic(ctx, svcCtx)
		resp, err := l.CreateSubscribeApplication(&input.Body)
		if err != nil {
			return nil, err
		}
		return &CreateSubscribeApplicationOutput{Body: resp}, nil
	}
}
