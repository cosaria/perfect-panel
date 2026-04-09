package subscribe

import (
	"context"

	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/models/subscribe"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/pkg/errors"
)

type GetSubscribeGroupListOutput struct {
	Body *types.GetSubscribeGroupListResponse
}

func GetSubscribeGroupListHandler(deps Deps) func(context.Context, *struct{}) (*GetSubscribeGroupListOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetSubscribeGroupListOutput, error) {
		l := NewGetSubscribeGroupListLogic(ctx, deps)
		resp, err := l.GetSubscribeGroupList()
		if err != nil {
			return nil, err
		}
		return &GetSubscribeGroupListOutput{Body: resp}, nil
	}
}

type GetSubscribeGroupListLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Get subscribe group list
func NewGetSubscribeGroupListLogic(ctx context.Context, deps Deps) *GetSubscribeGroupListLogic {
	return &GetSubscribeGroupListLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *GetSubscribeGroupListLogic) GetSubscribeGroupList() (resp *types.GetSubscribeGroupListResponse, err error) {
	var list []*subscribe.Group
	var total int64
	err = l.deps.DB.Model(&subscribe.Group{}).Count(&total).Find(&list).Error
	if err != nil {
		l.Error("[GetSubscribeGroupListLogic] get subscribe group list failed: ", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "get subscribe group list failed: %v", err.Error())
	}
	groupList := make([]types.SubscribeGroup, 0)
	tool.DeepCopy(&groupList, list)
	return &types.GetSubscribeGroupListResponse{
		Total: total,
		List:  groupList,
	}, nil
}
