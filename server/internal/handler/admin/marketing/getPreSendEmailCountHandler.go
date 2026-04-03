// huma:migrated
package marketing

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/marketing"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type GetPreSendEmailCountInput struct {
	Body types.GetPreSendEmailCountRequest
}

type GetPreSendEmailCountOutput struct {
	Body *types.GetPreSendEmailCountResponse
}

func GetPreSendEmailCountHandler(svcCtx *svc.ServiceContext) func(context.Context, *GetPreSendEmailCountInput) (*GetPreSendEmailCountOutput, error) {
	return func(ctx context.Context, input *GetPreSendEmailCountInput) (*GetPreSendEmailCountOutput, error) {
		l := marketing.NewGetPreSendEmailCountLogic(ctx, svcCtx)
		resp, err := l.GetPreSendEmailCount(&input.Body)
		if err != nil {
			return nil, err
		}
		return &GetPreSendEmailCountOutput{Body: resp}, nil
	}
}
