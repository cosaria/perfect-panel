// huma:migrated
package user

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/public/user"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type UpdateUserSubscribeNoteInput struct {
	Body types.UpdateUserSubscribeNoteRequest
}

func UpdateUserSubscribeNoteHandler(svcCtx *svc.ServiceContext) func(context.Context, *UpdateUserSubscribeNoteInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateUserSubscribeNoteInput) (*struct{}, error) {
		l := user.NewUpdateUserSubscribeNoteLogic(ctx, svcCtx)
		if err := l.UpdateUserSubscribeNote(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
