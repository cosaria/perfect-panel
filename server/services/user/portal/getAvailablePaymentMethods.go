package portal

import (
	"context"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type GetAvailablePaymentMethodsOutput struct {
	Body *types.GetAvailablePaymentMethodsResponse
}

func GetAvailablePaymentMethodsHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*GetAvailablePaymentMethodsOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetAvailablePaymentMethodsOutput, error) {
		l := NewGetAvailablePaymentMethodsLogic(ctx, svcCtx)
		resp, err := l.GetAvailablePaymentMethods()
		if err != nil {
			return nil, err
		}
		return &GetAvailablePaymentMethodsOutput{Body: resp}, nil
	}
}

type GetAvailablePaymentMethodsLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAvailablePaymentMethodsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAvailablePaymentMethodsLogic {
	return &GetAvailablePaymentMethodsLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAvailablePaymentMethodsLogic) GetAvailablePaymentMethods() (resp *types.GetAvailablePaymentMethodsResponse, err error) {
	data, err := l.svcCtx.PaymentModel.FindAvailableMethods(l.ctx)
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
