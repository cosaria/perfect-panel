// huma:migrated
package system

import (
	"context"
	"github.com/perfect-panel/server/services/admin/system"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type GetCurrencyConfigOutput struct {
	Body *types.CurrencyConfig
}

func GetCurrencyConfigHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*GetCurrencyConfigOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetCurrencyConfigOutput, error) {
		l := system.NewGetCurrencyConfigLogic(ctx, svcCtx)
		resp, err := l.GetCurrencyConfig()
		if err != nil {
			return nil, err
		}
		return &GetCurrencyConfigOutput{Body: resp}, nil
	}
}
