package common

import (
	"context"
	"encoding/json"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/pkg/errors"
)

type GetClientOutput struct {
	Body *types.GetSubscribeClientResponse
}

func GetClientHandler(deps Deps) func(context.Context, *struct{}) (*GetClientOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetClientOutput, error) {
		l := NewGetClientLogic(ctx, deps)
		resp, err := l.GetClient()
		if err != nil {
			return nil, err
		}
		return &GetClientOutput{Body: resp}, nil
	}
}

type GetClientLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Get Client
func NewGetClientLogic(ctx context.Context, deps Deps) *GetClientLogic {
	return &GetClientLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *GetClientLogic) GetClient() (resp *types.GetSubscribeClientResponse, err error) {
	data, err := l.deps.ClientModel.List(l.ctx)
	if err != nil {
		l.Errorf("Failed to get subscribe application list: %v", err)
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "Failed to get subscribe application list")
	}
	var list []types.SubscribeClient
	for _, item := range data {
		var temp types.DownloadLink
		if item.DownloadLink != "" {
			_ = json.Unmarshal([]byte(item.DownloadLink), &temp)
		}
		list = append(list, types.SubscribeClient{
			Id:           item.Id,
			Name:         item.Name,
			Description:  item.Description,
			Icon:         item.Icon,
			Scheme:       item.Scheme,
			IsDefault:    item.IsDefault,
			DownloadLink: temp,
		})
	}
	resp = &types.GetSubscribeClientResponse{
		Total: int64(len(list)),
		List:  list,
	}
	return
}
