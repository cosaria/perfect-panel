package authMethod

import (
	"context"
	"encoding/json"

	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/pkg/errors"
)

type GetAuthMethodListOutput struct {
	Body *types.GetAuthMethodListResponse
}

func GetAuthMethodListHandler(deps Deps) func(context.Context, *struct{}) (*GetAuthMethodListOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetAuthMethodListOutput, error) {
		l := NewGetAuthMethodListLogic(ctx, deps)
		resp, err := l.GetAuthMethodList()
		if err != nil {
			return nil, err
		}
		return &GetAuthMethodListOutput{Body: resp}, nil
	}
}

type GetAuthMethodListLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// NewGetAuthMethodListLogic Get auth method list
func NewGetAuthMethodListLogic(ctx context.Context, deps Deps) *GetAuthMethodListLogic {
	return &GetAuthMethodListLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *GetAuthMethodListLogic) GetAuthMethodList() (resp *types.GetAuthMethodListResponse, err error) {
	methods, err := l.deps.AuthModel.FindAll(l.ctx)
	if err != nil {
		l.Errorw("find all failed", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "find all failed: %v", err.Error())
	}
	var list []types.AuthMethodConfig
	for _, method := range methods {
		var item types.AuthMethodConfig
		tool.DeepCopy(&item, method)
		if method.Config != "" {
			if err := json.Unmarshal([]byte(method.Config), &item.Config); err != nil {
				l.Errorw("unmarshal config failed", logger.Field("config", method.Config), logger.Field("error", err.Error()))
				return nil, errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "unmarshal config failed: %v", err.Error())
			}
		}
		list = append(list, item)
	}
	return &types.GetAuthMethodListResponse{List: list}, nil
}
