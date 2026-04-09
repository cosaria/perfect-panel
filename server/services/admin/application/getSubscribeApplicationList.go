package application

import (
	"context"
	"encoding/json"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/pkg/errors"
)

type GetSubscribeApplicationListInput struct {
	types.GetSubscribeApplicationListRequest
}

type GetSubscribeApplicationListOutput struct {
	Body *types.GetSubscribeApplicationListResponse
}

func GetSubscribeApplicationListHandler(deps Deps) func(context.Context, *GetSubscribeApplicationListInput) (*GetSubscribeApplicationListOutput, error) {
	return func(ctx context.Context, input *GetSubscribeApplicationListInput) (*GetSubscribeApplicationListOutput, error) {
		l := NewGetSubscribeApplicationListLogic(ctx, deps)
		resp, err := l.GetSubscribeApplicationList(&input.GetSubscribeApplicationListRequest)
		if err != nil {
			return nil, err
		}
		return &GetSubscribeApplicationListOutput{Body: resp}, nil
	}
}

type GetSubscribeApplicationListLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// NewGetSubscribeApplicationListLogic Get subscribe application list
func NewGetSubscribeApplicationListLogic(ctx context.Context, deps Deps) *GetSubscribeApplicationListLogic {
	return &GetSubscribeApplicationListLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *GetSubscribeApplicationListLogic) GetSubscribeApplicationList(req *types.GetSubscribeApplicationListRequest) (resp *types.GetSubscribeApplicationListResponse, err error) {
	data, err := l.deps.ClientModel.List(l.ctx)
	if err != nil {
		l.Errorf("Failed to get subscribe application list: %v", err)
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "Failed to get subscribe application list")
	}
	var list []types.SubscribeApplication
	for _, item := range data {
		var temp types.DownloadLink
		if item.DownloadLink != "" {
			_ = json.Unmarshal([]byte(item.DownloadLink), &temp)
		}
		list = append(list, types.SubscribeApplication{
			Id:                item.Id,
			Name:              item.Name,
			Description:       item.Description,
			Icon:              item.Icon,
			Scheme:            item.Scheme,
			UserAgent:         item.UserAgent,
			IsDefault:         item.IsDefault,
			SubscribeTemplate: item.SubscribeTemplate,
			OutputFormat:      item.OutputFormat,
			DownloadLink:      temp,
			CreatedAt:         item.CreatedAt.UnixMilli(),
			UpdatedAt:         item.UpdatedAt.UnixMilli(),
		})
	}
	resp = &types.GetSubscribeApplicationListResponse{
		Total: int64(len(list)),
		List:  list,
	}
	return
}
