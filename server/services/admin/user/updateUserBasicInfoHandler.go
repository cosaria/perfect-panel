// huma:migrated
package user

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type UpdateUserBasicInfoInput struct {
	Body types.UpdateUserBasiceInfoRequest
}

func UpdateUserBasicInfoHandler(svcCtx *svc.ServiceContext) func(context.Context, *UpdateUserBasicInfoInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateUserBasicInfoInput) (*struct{}, error) {
		l := NewUpdateUserBasicInfoLogic(ctx, svcCtx)
		if err := l.UpdateUserBasicInfo(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
