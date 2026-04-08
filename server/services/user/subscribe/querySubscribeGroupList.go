package subscribe

import (
	"context"
	"github.com/perfect-panel/server/models/subscribe"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type QuerySubscribeGroupListOutput struct {
	Body *types.QuerySubscribeGroupListResponse
}

func QuerySubscribeGroupListHandler(deps Deps) func(context.Context, *struct{}) (*QuerySubscribeGroupListOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*QuerySubscribeGroupListOutput, error) {
		l := NewQuerySubscribeGroupListLogic(ctx, deps)
		resp, err := l.QuerySubscribeGroupList()
		if err != nil {
			return nil, err
		}
		return &QuerySubscribeGroupListOutput{Body: resp}, nil
	}
}

type QuerySubscribeGroupListLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Get subscribe group list
func NewQuerySubscribeGroupListLogic(ctx context.Context, deps Deps) *QuerySubscribeGroupListLogic {
	return &QuerySubscribeGroupListLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *QuerySubscribeGroupListLogic) QuerySubscribeGroupList() (resp *types.QuerySubscribeGroupListResponse, err error) {
	var list []*subscribe.Group
	var total int64
	err = l.deps.DB.Model(&subscribe.Group{}).Count(&total).Find(&list).Error
	if err != nil {
		l.Error("[QuerySubscribeGroupListLogic] get subscribe group list failed: ", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "get subscribe group list failed: %v", err.Error())
	}
	groupList := make([]types.SubscribeGroup, 0)
	tool.DeepCopy(&groupList, list)
	return &types.QuerySubscribeGroupListResponse{
		Total: total,
		List:  groupList,
	}, nil
}
