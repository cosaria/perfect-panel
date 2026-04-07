// huma:migrated
package user

import (
	"context"
	"github.com/perfect-panel/server/services/user/user"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type UpdateUserRulesInput struct {
	Body types.UpdateUserRulesRequest
}

func UpdateUserRulesHandler(svcCtx *svc.ServiceContext) func(context.Context, *UpdateUserRulesInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateUserRulesInput) (*struct{}, error) {
		l := user.NewUpdateUserRulesLogic(ctx, svcCtx)
		if err := l.UpdateUserRules(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
