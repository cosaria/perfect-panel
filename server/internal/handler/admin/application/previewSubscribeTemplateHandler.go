// huma:migrated
package application

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/application"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type PreviewSubscribeTemplateInput struct {
	types.PreviewSubscribeTemplateRequest
}

type PreviewSubscribeTemplateOutput struct {
	Body *types.PreviewSubscribeTemplateResponse
}

func PreviewSubscribeTemplateHandler(svcCtx *svc.ServiceContext) func(context.Context, *PreviewSubscribeTemplateInput) (*PreviewSubscribeTemplateOutput, error) {
	return func(ctx context.Context, input *PreviewSubscribeTemplateInput) (*PreviewSubscribeTemplateOutput, error) {
		l := application.NewPreviewSubscribeTemplateLogic(ctx, svcCtx)
		resp, err := l.PreviewSubscribeTemplate(&input.PreviewSubscribeTemplateRequest)
		if err != nil {
			return nil, err
		}
		return &PreviewSubscribeTemplateOutput{Body: resp}, nil
	}
}
