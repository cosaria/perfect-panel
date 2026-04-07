// huma:migrated
package payment

import (
	"context"
	"github.com/perfect-panel/server/services/admin/payment"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type GetPaymentMethodListInput struct {
	Body types.GetPaymentMethodListRequest
}

type GetPaymentMethodListOutput struct {
	Body *types.GetPaymentMethodListResponse
}

func GetPaymentMethodListHandler(svcCtx *svc.ServiceContext) func(context.Context, *GetPaymentMethodListInput) (*GetPaymentMethodListOutput, error) {
	return func(ctx context.Context, input *GetPaymentMethodListInput) (*GetPaymentMethodListOutput, error) {
		l := payment.NewGetPaymentMethodListLogic(ctx, svcCtx)
		resp, err := l.GetPaymentMethodList(&input.Body)
		if err != nil {
			return nil, err
		}
		return &GetPaymentMethodListOutput{Body: resp}, nil
	}
}
