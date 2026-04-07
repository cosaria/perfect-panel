// huma:migrated
package user

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type UpdateBindMobileInput struct {
	Body types.UpdateBindMobileRequest
}

func UpdateBindMobileHandler(svcCtx *svc.ServiceContext) func(context.Context, *UpdateBindMobileInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateBindMobileInput) (*struct{}, error) {
		l := NewUpdateBindMobileLogic(ctx, svcCtx)
		if err := l.UpdateBindMobile(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
