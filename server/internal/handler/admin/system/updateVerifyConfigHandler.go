// huma:migrated
package system

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/system"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type UpdateVerifyConfigInput struct {
	Body types.VerifyConfig
}

func UpdateVerifyConfigHandler(svcCtx *svc.ServiceContext) func(context.Context, *UpdateVerifyConfigInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateVerifyConfigInput) (*struct{}, error) {
		l := system.NewUpdateVerifyConfigLogic(ctx, svcCtx)
		if err := l.UpdateVerifyConfig(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
