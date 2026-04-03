// huma:migrated
package system

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/system"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type UpdateCurrencyConfigInput struct {
	Body types.CurrencyConfig
}

func UpdateCurrencyConfigHandler(svcCtx *svc.ServiceContext) func(context.Context, *UpdateCurrencyConfigInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateCurrencyConfigInput) (*struct{}, error) {
		l := system.NewUpdateCurrencyConfigLogic(ctx, svcCtx)
		if err := l.UpdateCurrencyConfig(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
