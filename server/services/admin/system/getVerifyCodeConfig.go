package system

import (
	"context"

	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type GetVerifyCodeConfigOutput struct {
	Body *types.VerifyCodeConfig
}

func GetVerifyCodeConfigHandler(deps Deps) func(context.Context, *struct{}) (*GetVerifyCodeConfigOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetVerifyCodeConfigOutput, error) {
		l := NewGetVerifyCodeConfigLogic(ctx, deps)
		resp, err := l.GetVerifyCodeConfig()
		if err != nil {
			return nil, err
		}
		return &GetVerifyCodeConfigOutput{Body: resp}, nil
	}
}

type GetVerifyCodeConfigLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Get Verify Code Config
func NewGetVerifyCodeConfigLogic(ctx context.Context, deps Deps) *GetVerifyCodeConfigLogic {
	return &GetVerifyCodeConfigLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *GetVerifyCodeConfigLogic) GetVerifyCodeConfig() (resp *types.VerifyCodeConfig, err error) {
	data, err := l.deps.SystemModel.GetVerifyCodeConfig(l.ctx)
	if err != nil {
		l.Errorw("Get Verify Code Config Error: ", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "Get Verify Code Config Error: %s", err.Error())
	}
	resp = &types.VerifyCodeConfig{}
	tool.SystemConfigSliceReflectToStruct(data, resp)
	return
}
