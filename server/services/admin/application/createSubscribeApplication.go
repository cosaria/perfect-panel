package application

import (
	"context"
	"github.com/perfect-panel/server/models/client"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type CreateSubscribeApplicationInput struct {
	Body types.CreateSubscribeApplicationRequest
}

type CreateSubscribeApplicationOutput struct {
	Body *types.SubscribeApplication
}

func CreateSubscribeApplicationHandler(svcCtx *svc.ServiceContext) func(context.Context, *CreateSubscribeApplicationInput) (*CreateSubscribeApplicationOutput, error) {
	return func(ctx context.Context, input *CreateSubscribeApplicationInput) (*CreateSubscribeApplicationOutput, error) {
		l := NewCreateSubscribeApplicationLogic(ctx, svcCtx)
		resp, err := l.CreateSubscribeApplication(&input.Body)
		if err != nil {
			return nil, err
		}
		return &CreateSubscribeApplicationOutput{Body: resp}, nil
	}
}

type CreateSubscribeApplicationLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewCreateSubscribeApplicationLogic Create subscribe application
func NewCreateSubscribeApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateSubscribeApplicationLogic {
	return &CreateSubscribeApplicationLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateSubscribeApplicationLogic) CreateSubscribeApplication(req *types.CreateSubscribeApplicationRequest) (resp *types.SubscribeApplication, err error) {
	var link client.DownloadLink
	tool.DeepCopy(&link, req.DownloadLink)
	linkData, err := link.Marshal()
	if err != nil {
		l.Errorf("Failed to marshal download link: %v", err)
		return nil, errors.Wrap(xerr.NewErrCode(xerr.ERROR), " Failed to marshal download link")
	}
	data := &client.SubscribeApplication{
		Name:              req.Name,
		Icon:              req.Icon,
		Description:       req.Description,
		Scheme:            req.Scheme,
		UserAgent:         req.UserAgent,
		IsDefault:         req.IsDefault,
		SubscribeTemplate: req.SubscribeTemplate,
		OutputFormat:      req.OutputFormat,
		DownloadLink:      string(linkData),
	}

	err = l.svcCtx.ClientModel.Insert(l.ctx, data)
	if err != nil {
		l.Errorf("Failed to create subscribe application: %v", err)
		return nil, errors.Wrap(xerr.NewErrCode(xerr.DatabaseInsertError), "Failed to create subscribe application")
	}

	resp = &types.SubscribeApplication{}
	tool.DeepCopy(resp, data)
	resp.DownloadLink = req.DownloadLink

	return
}
