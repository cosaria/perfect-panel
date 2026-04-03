// huma:migrated
package authMethod

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/authMethod"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type GetAuthMethodListOutput struct {
	Body *types.GetAuthMethodListResponse
}

func GetAuthMethodListHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*GetAuthMethodListOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetAuthMethodListOutput, error) {
		l := authMethod.NewGetAuthMethodListLogic(ctx, svcCtx)
		resp, err := l.GetAuthMethodList()
		if err != nil {
			return nil, err
		}
		return &GetAuthMethodListOutput{Body: resp}, nil
	}
}
