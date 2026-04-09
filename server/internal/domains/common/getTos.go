package common

import (
	"context"

	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/pkg/errors"
)

type GetTosOutput struct {
	Body *types.GetTosResponse
}

func GetTosHandler(deps Deps) func(context.Context, *struct{}) (*GetTosOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetTosOutput, error) {
		l := NewGetTosLogic(ctx, deps)
		resp, err := l.GetTos()
		if err != nil {
			return nil, err
		}
		return &GetTosOutput{Body: resp}, nil
	}
}

type GetTosLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Get Tos
func NewGetTosLogic(ctx context.Context, deps Deps) *GetTosLogic {
	return &GetTosLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *GetTosLogic) GetTos() (resp *types.GetTosResponse, err error) {
	resp = &types.GetTosResponse{}
	// get Tos config from db
	configs, err := l.deps.SystemModel.GetTosConfig(l.ctx)
	if err != nil {
		l.Errorw("[GetTosLogic] GetTos error: ", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "GetTos error: %v", err.Error())
	}
	// reflect to response
	tool.SystemConfigSliceReflectToStruct(configs, resp)
	return
}
