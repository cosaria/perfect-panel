package system

import (
	"context"

	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/pkg/errors"
)

type GetPrivacyPolicyConfigOutput struct {
	Body *types.PrivacyPolicyConfig
}

func GetPrivacyPolicyConfigHandler(deps Deps) func(context.Context, *struct{}) (*GetPrivacyPolicyConfigOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetPrivacyPolicyConfigOutput, error) {
		l := NewGetPrivacyPolicyConfigLogic(ctx, deps)
		resp, err := l.GetPrivacyPolicyConfig()
		if err != nil {
			return nil, err
		}
		return &GetPrivacyPolicyConfigOutput{Body: resp}, nil
	}
}

type GetPrivacyPolicyConfigLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// NewGetPrivacyPolicyConfigLogic get Privacy Policy Config
func NewGetPrivacyPolicyConfigLogic(ctx context.Context, deps Deps) *GetPrivacyPolicyConfigLogic {
	return &GetPrivacyPolicyConfigLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *GetPrivacyPolicyConfigLogic) GetPrivacyPolicyConfig() (resp *types.PrivacyPolicyConfig, err error) {
	resp = &types.PrivacyPolicyConfig{}
	// get tos config from db
	configs, err := l.deps.SystemModel.GetTosConfig(l.ctx)
	if err != nil {
		l.Errorw("[GetTosConfig] GetTosConfig error", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "GetTosConfig error: %v", err.Error())
	}
	// reflect to response
	tool.SystemConfigSliceReflectToStruct(configs, resp)
	return
}
