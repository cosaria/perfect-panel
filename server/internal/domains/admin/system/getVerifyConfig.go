package system

import (
	"context"

	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/pkg/errors"
)

type GetVerifyConfigOutput struct {
	Body *types.VerifyConfig
}

func GetVerifyConfigHandler(deps Deps) func(context.Context, *struct{}) (*GetVerifyConfigOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetVerifyConfigOutput, error) {
		l := NewGetVerifyConfigLogic(ctx, deps)
		resp, err := l.GetVerifyConfig()
		if err != nil {
			return nil, err
		}
		return &GetVerifyConfigOutput{Body: resp}, nil
	}
}

type GetVerifyConfigLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

func NewGetVerifyConfigLogic(ctx context.Context, deps Deps) *GetVerifyConfigLogic {
	return &GetVerifyConfigLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *GetVerifyConfigLogic) GetVerifyConfig() (*types.VerifyConfig, error) {
	resp := &types.VerifyConfig{}
	// get verify config from db
	verifyConfigs, err := l.deps.SystemModel.GetVerifyConfig(l.ctx)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "get verify config failed: %v", err.Error())
	}
	// reflect to response
	tool.SystemConfigSliceReflectToStruct(verifyConfigs, resp)
	// update verify config to system
	if l.deps.ReloadVerify != nil {
		l.deps.ReloadVerify()
	}
	return resp, nil
}
