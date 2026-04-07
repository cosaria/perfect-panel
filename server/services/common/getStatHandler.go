// huma:migrated
package common

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type GetStatOutput struct {
	Body *types.GetStatResponse
}

func GetStatHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*GetStatOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetStatOutput, error) {
		l := NewGetStatLogic(ctx, svcCtx)
		resp, err := l.GetStat()
		if err != nil {
			return nil, err
		}
		return &GetStatOutput{Body: resp}, nil
	}
}
