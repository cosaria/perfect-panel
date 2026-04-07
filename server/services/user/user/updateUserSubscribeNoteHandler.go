// huma:migrated
package user

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type UpdateUserSubscribeNoteInput struct {
	Body types.UpdateUserSubscribeNoteRequest
}

func UpdateUserSubscribeNoteHandler(svcCtx *svc.ServiceContext) func(context.Context, *UpdateUserSubscribeNoteInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateUserSubscribeNoteInput) (*struct{}, error) {
		l := NewUpdateUserSubscribeNoteLogic(ctx, svcCtx)
		if err := l.UpdateUserSubscribeNote(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
