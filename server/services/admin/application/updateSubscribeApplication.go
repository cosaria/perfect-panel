package application

import (
	"context"
	"github.com/perfect-panel/server/models/client"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type UpdateSubscribeApplicationInput struct {
	Body types.UpdateSubscribeApplicationRequest
}

type UpdateSubscribeApplicationOutput struct {
	Body *types.SubscribeApplication
}

func UpdateSubscribeApplicationHandler(deps Deps) func(context.Context, *UpdateSubscribeApplicationInput) (*UpdateSubscribeApplicationOutput, error) {
	return func(ctx context.Context, input *UpdateSubscribeApplicationInput) (*UpdateSubscribeApplicationOutput, error) {
		l := NewUpdateSubscribeApplicationLogic(ctx, deps)
		resp, err := l.UpdateSubscribeApplication(&input.Body)
		if err != nil {
			return nil, err
		}
		return &UpdateSubscribeApplicationOutput{Body: resp}, nil
	}
}

type UpdateSubscribeApplicationLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// NewUpdateSubscribeApplicationLogic Update subscribe application
func NewUpdateSubscribeApplicationLogic(ctx context.Context, deps Deps) *UpdateSubscribeApplicationLogic {
	return &UpdateSubscribeApplicationLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *UpdateSubscribeApplicationLogic) UpdateSubscribeApplication(req *types.UpdateSubscribeApplicationRequest) (resp *types.SubscribeApplication, err error) {
	data, err := l.deps.ClientModel.FindOne(l.ctx, req.Id)
	if err != nil {
		l.Errorf("Failed to find subscribe application with ID %d: %v", req.Id, err)
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "Failed to find subscribe application with ID %d", req.Id)
	}
	var link client.DownloadLink
	tool.DeepCopy(&link, req.DownloadLink)
	linkData, err := link.Marshal()
	if err != nil {
		l.Errorf("Failed to marshal download link: %v", err)
		return nil, errors.Wrap(xerr.NewErrCode(xerr.ERROR), " Failed to marshal download link")
	}

	data.Name = req.Name
	data.Icon = req.Icon
	data.Description = req.Description
	data.Scheme = req.Scheme
	data.UserAgent = req.UserAgent
	data.IsDefault = req.IsDefault
	data.SubscribeTemplate = req.SubscribeTemplate
	data.OutputFormat = req.OutputFormat
	data.DownloadLink = string(linkData)
	err = l.deps.ClientModel.Update(l.ctx, data)
	if err != nil {
		l.Errorf("Failed to update subscribe application with ID %d: %v", req.Id, err)
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseUpdateError), "Failed to update subscribe application with ID %d", req.Id)
	}
	resp = &types.SubscribeApplication{}
	tool.DeepCopy(&resp, data)
	resp.DownloadLink = req.DownloadLink
	return
}
