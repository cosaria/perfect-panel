// huma:migrated
package system

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/system"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type UpdatePrivacyPolicyConfigInput struct {
	Body types.PrivacyPolicyConfig
}

func UpdatePrivacyPolicyConfigHandler(svcCtx *svc.ServiceContext) func(context.Context, *UpdatePrivacyPolicyConfigInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdatePrivacyPolicyConfigInput) (*struct{}, error) {
		l := system.NewUpdatePrivacyPolicyConfigLogic(ctx, svcCtx)
		if err := l.UpdatePrivacyPolicyConfig(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
