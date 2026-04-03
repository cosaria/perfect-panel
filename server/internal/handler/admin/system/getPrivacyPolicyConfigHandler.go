// huma:migrated
package system

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/system"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type GetPrivacyPolicyConfigOutput struct {
	Body *types.PrivacyPolicyConfig
}

func GetPrivacyPolicyConfigHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*GetPrivacyPolicyConfigOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetPrivacyPolicyConfigOutput, error) {
		l := system.NewGetPrivacyPolicyConfigLogic(ctx, svcCtx)
		resp, err := l.GetPrivacyPolicyConfig()
		if err != nil {
			return nil, err
		}
		return &GetPrivacyPolicyConfigOutput{Body: resp}, nil
	}
}
