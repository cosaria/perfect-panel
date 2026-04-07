package common

import (
	"context"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type GetPrivacyPolicyOutput struct {
	Body *types.PrivacyPolicyConfig
}

func GetPrivacyPolicyHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*GetPrivacyPolicyOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetPrivacyPolicyOutput, error) {
		l := NewGetPrivacyPolicyLogic(ctx, svcCtx)
		resp, err := l.GetPrivacyPolicy()
		if err != nil {
			return nil, err
		}
		return &GetPrivacyPolicyOutput{Body: resp}, nil
	}
}

type GetPrivacyPolicyLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Get Privacy Policy
func NewGetPrivacyPolicyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPrivacyPolicyLogic {
	return &GetPrivacyPolicyLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPrivacyPolicyLogic) GetPrivacyPolicy() (resp *types.PrivacyPolicyConfig, err error) {
	resp = &types.PrivacyPolicyConfig{}
	// get tos config from db
	configs, err := l.svcCtx.SystemModel.GetTosConfig(l.ctx)
	if err != nil {
		l.Errorw("[GetTosConfig] GetTosConfig error", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "GetTosConfig error: %v", err.Error())
	}
	// reflect to response
	tool.SystemConfigSliceReflectToStruct(configs, resp)
	return
}
