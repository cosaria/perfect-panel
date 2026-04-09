package application

import (
	"context"
	"time"

	"github.com/perfect-panel/server/adapter"
	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/models/node"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/pkg/errors"
)

type PreviewSubscribeTemplateInput struct {
	types.PreviewSubscribeTemplateRequest
}

type PreviewSubscribeTemplateOutput struct {
	Body *types.PreviewSubscribeTemplateResponse
}

func PreviewSubscribeTemplateHandler(deps Deps) func(context.Context, *PreviewSubscribeTemplateInput) (*PreviewSubscribeTemplateOutput, error) {
	return func(ctx context.Context, input *PreviewSubscribeTemplateInput) (*PreviewSubscribeTemplateOutput, error) {
		l := NewPreviewSubscribeTemplateLogic(ctx, deps)
		resp, err := l.PreviewSubscribeTemplate(&input.PreviewSubscribeTemplateRequest)
		if err != nil {
			return nil, err
		}
		return &PreviewSubscribeTemplateOutput{Body: resp}, nil
	}
}

type PreviewSubscribeTemplateLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Preview Template
func NewPreviewSubscribeTemplateLogic(ctx context.Context, deps Deps) *PreviewSubscribeTemplateLogic {
	return &PreviewSubscribeTemplateLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *PreviewSubscribeTemplateLogic) PreviewSubscribeTemplate(req *types.PreviewSubscribeTemplateRequest) (resp *types.PreviewSubscribeTemplateResponse, err error) {
	enable := true
	_, servers, err := l.deps.NodeModel.FilterNodeList(l.ctx, &node.FilterNodeParams{
		Page:    1,
		Size:    1000,
		Preload: true,
		Enabled: &enable,
	})
	if err != nil {
		l.Errorf("[PreviewSubscribeTemplateLogic] FindAllServer error: %v", err.Error())
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "FindAllServer error: %v", err.Error())
	}

	data, err := l.deps.ClientModel.FindOne(l.ctx, req.Id)
	if err != nil {
		l.Errorf("[PreviewSubscribeTemplateLogic] FindOne error: %v", err.Error())
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "FindOneClient error: %v", err.Error())
	}

	sub := adapter.NewAdapter(data.SubscribeTemplate, adapter.WithServers(servers),
		adapter.WithSiteName("PerfectPanel"),
		adapter.WithSubscribeName("Test Subscribe"),
		adapter.WithOutputFormat(data.OutputFormat),
		adapter.WithUserInfo(adapter.User{
			Password:     "test-password",
			ExpiredAt:    time.Now().AddDate(1, 0, 0),
			Download:     0,
			Upload:       0,
			Traffic:      1000,
			SubscribeURL: "https://example.com/subscribe",
		}))
	// Get client config
	a, err := sub.Client()
	if err != nil {
		l.Errorf("[PreviewSubscribeTemplateLogic] Client error: %v", err.Error())
		return nil, errors.Wrapf(xerr.NewErrMsg(err.Error()), "Client error: %v", err.Error())
	}
	bytes, err := a.Build()
	if err != nil {
		l.Errorf("[PreviewSubscribeTemplateLogic] Build error: %v", err.Error())
		return nil, errors.Wrapf(xerr.NewErrMsg(err.Error()), "Build error: %v", err.Error())
	}
	return &types.PreviewSubscribeTemplateResponse{
		Template: string(bytes),
	}, nil
}
