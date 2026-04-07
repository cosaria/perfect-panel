// huma:migrated
package system

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type GetPrivacyPolicyConfigOutput struct {
	Body *types.PrivacyPolicyConfig
}

func GetPrivacyPolicyConfigHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*GetPrivacyPolicyConfigOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetPrivacyPolicyConfigOutput, error) {
		l := NewGetPrivacyPolicyConfigLogic(ctx, svcCtx)
		resp, err := l.GetPrivacyPolicyConfig()
		if err != nil {
			return nil, err
		}
		return &GetPrivacyPolicyConfigOutput{Body: resp}, nil
	}
}
