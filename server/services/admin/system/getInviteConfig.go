package system

import (
	"context"

	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type GetInviteConfigOutput struct {
	Body *types.InviteConfig
}

func GetInviteConfigHandler(deps Deps) func(context.Context, *struct{}) (*GetInviteConfigOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetInviteConfigOutput, error) {
		l := NewGetInviteConfigLogic(ctx, deps)
		resp, err := l.GetInviteConfig()
		if err != nil {
			return nil, err
		}
		return &GetInviteConfigOutput{Body: resp}, nil
	}
}

type GetInviteConfigLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

func NewGetInviteConfigLogic(ctx context.Context, deps Deps) *GetInviteConfigLogic {
	return &GetInviteConfigLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *GetInviteConfigLogic) GetInviteConfig() (*types.InviteConfig, error) {
	resp := &types.InviteConfig{}
	// get invite config from db
	configs, err := l.deps.SystemModel.GetInviteConfig(l.ctx)
	if err != nil {
		l.Errorw("[GetInviteConfigLogic] get invite config error: ", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "get invite config error: %v", err.Error())
	}
	// reflect to response
	tool.SystemConfigSliceReflectToStruct(configs, resp)

	return resp, nil
}
