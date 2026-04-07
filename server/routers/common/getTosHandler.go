// huma:migrated
package common

import (
	"context"
	"github.com/perfect-panel/server/services/common"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type GetTosOutput struct {
	Body *types.GetTosResponse
}

func GetTosHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*GetTosOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetTosOutput, error) {
		l := common.NewGetTosLogic(ctx, svcCtx)
		resp, err := l.GetTos()
		if err != nil {
			return nil, err
		}
		return &GetTosOutput{Body: resp}, nil
	}
}
