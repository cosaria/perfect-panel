// huma:migrated
package user

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type GetUserDetailInput struct {
	types.GetDetailRequest
}

type GetUserDetailOutput struct {
	Body *types.User
}

func GetUserDetailHandler(svcCtx *svc.ServiceContext) func(context.Context, *GetUserDetailInput) (*GetUserDetailOutput, error) {
	return func(ctx context.Context, input *GetUserDetailInput) (*GetUserDetailOutput, error) {
		l := NewGetUserDetailLogic(ctx, svcCtx)
		resp, err := l.GetUserDetail(&input.GetDetailRequest)
		if err != nil {
			return nil, err
		}
		return &GetUserDetailOutput{Body: resp}, nil
	}
}
