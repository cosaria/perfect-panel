// huma:migrated
package portal

import (
	"context"
	"github.com/perfect-panel/server/services/user/portal"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type GetAvailablePaymentMethodsOutput struct {
	Body *types.GetAvailablePaymentMethodsResponse
}

func GetAvailablePaymentMethodsHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*GetAvailablePaymentMethodsOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetAvailablePaymentMethodsOutput, error) {
		l := portal.NewGetAvailablePaymentMethodsLogic(ctx, svcCtx)
		resp, err := l.GetAvailablePaymentMethods()
		if err != nil {
			return nil, err
		}
		return &GetAvailablePaymentMethodsOutput{Body: resp}, nil
	}
}
