// huma:migrated
package user

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/public/user"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type UpdateBindMobileInput struct {
	Body types.UpdateBindMobileRequest
}

func UpdateBindMobileHandler(svcCtx *svc.ServiceContext) func(context.Context, *UpdateBindMobileInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateBindMobileInput) (*struct{}, error) {
		l := user.NewUpdateBindMobileLogic(ctx, svcCtx)
		if err := l.UpdateBindMobile(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
