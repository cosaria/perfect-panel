// huma:migrated
package system

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/system"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type UpdateVerifyCodeConfigInput struct {
	Body types.VerifyCodeConfig
}

func UpdateVerifyCodeConfigHandler(svcCtx *svc.ServiceContext) func(context.Context, *UpdateVerifyCodeConfigInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateVerifyCodeConfigInput) (*struct{}, error) {
		l := system.NewUpdateVerifyCodeConfigLogic(ctx, svcCtx)
		if err := l.UpdateVerifyCodeConfig(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
