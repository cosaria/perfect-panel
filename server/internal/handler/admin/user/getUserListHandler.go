// huma:migrated
package user

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/user"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type GetUserListInput struct {
	Body types.GetUserListRequest
}

type GetUserListOutput struct {
	Body *types.GetUserListResponse
}

func GetUserListHandler(svcCtx *svc.ServiceContext) func(context.Context, *GetUserListInput) (*GetUserListOutput, error) {
	return func(ctx context.Context, input *GetUserListInput) (*GetUserListOutput, error) {
		l := user.NewGetUserListLogic(ctx, svcCtx)
		resp, err := l.GetUserList(&input.Body)
		if err != nil {
			return nil, err
		}
		return &GetUserListOutput{Body: resp}, nil
	}
}
