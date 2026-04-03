// huma:migrated
package common

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/common"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type GetPrivacyPolicyOutput struct {
	Body *types.PrivacyPolicyConfig
}

func GetPrivacyPolicyHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*GetPrivacyPolicyOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetPrivacyPolicyOutput, error) {
		l := common.NewGetPrivacyPolicyLogic(ctx, svcCtx)
		resp, err := l.GetPrivacyPolicy()
		if err != nil {
			return nil, err
		}
		return &GetPrivacyPolicyOutput{Body: resp}, nil
	}
}
