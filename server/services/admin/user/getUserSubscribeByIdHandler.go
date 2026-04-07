// huma:migrated
package user

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type GetUserSubscribeByIdInput struct {
	types.GetUserSubscribeByIdRequest
}

type GetUserSubscribeByIdOutput struct {
	Body *types.UserSubscribeDetail
}

func GetUserSubscribeByIdHandler(svcCtx *svc.ServiceContext) func(context.Context, *GetUserSubscribeByIdInput) (*GetUserSubscribeByIdOutput, error) {
	return func(ctx context.Context, input *GetUserSubscribeByIdInput) (*GetUserSubscribeByIdOutput, error) {
		l := NewGetUserSubscribeByIdLogic(ctx, svcCtx)
		resp, err := l.GetUserSubscribeById(&input.GetUserSubscribeByIdRequest)
		if err != nil {
			return nil, err
		}
		return &GetUserSubscribeByIdOutput{Body: resp}, nil
	}
}
