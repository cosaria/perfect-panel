package payment

import (
	"context"

	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/pkg/errors"
)

type GetAvailablePaymentMethodsOutput struct {
	Body *types.GetAvailablePaymentMethodsResponse
}

func GetAvailablePaymentMethodsHandler(deps Deps) func(context.Context, *struct{}) (*GetAvailablePaymentMethodsOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetAvailablePaymentMethodsOutput, error) {
		l := NewGetAvailablePaymentMethodsLogic(ctx, deps)
		resp, err := l.GetAvailablePaymentMethods()
		if err != nil {
			return nil, err
		}
		return &GetAvailablePaymentMethodsOutput{Body: resp}, nil
	}
}

type GetAvailablePaymentMethodsLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Get available payment methods
func NewGetAvailablePaymentMethodsLogic(ctx context.Context, deps Deps) *GetAvailablePaymentMethodsLogic {
	return &GetAvailablePaymentMethodsLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *GetAvailablePaymentMethodsLogic) GetAvailablePaymentMethods() (resp *types.GetAvailablePaymentMethodsResponse, err error) {
	data, err := l.deps.PaymentModel.FindAvailableMethods(l.ctx)
	if err != nil {
		l.Errorw("[GetAvailablePaymentMethods] database error", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "GetAvailablePaymentMethods: %v", err.Error())
	}
	resp = &types.GetAvailablePaymentMethodsResponse{
		List: make([]types.PaymentMethod, 0),
	}

	tool.DeepCopy(&resp.List, data)

	return
}
