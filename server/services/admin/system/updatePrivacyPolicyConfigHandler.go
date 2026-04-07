// huma:migrated
package system

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type UpdatePrivacyPolicyConfigInput struct {
	Body types.PrivacyPolicyConfig
}

func UpdatePrivacyPolicyConfigHandler(svcCtx *svc.ServiceContext) func(context.Context, *UpdatePrivacyPolicyConfigInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdatePrivacyPolicyConfigInput) (*struct{}, error) {
		l := NewUpdatePrivacyPolicyConfigLogic(ctx, svcCtx)
		if err := l.UpdatePrivacyPolicyConfig(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
