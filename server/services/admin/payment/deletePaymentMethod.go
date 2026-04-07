package payment

import (
	"context"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type DeletePaymentMethodInput struct {
	Body types.DeletePaymentMethodRequest
}

func DeletePaymentMethodHandler(svcCtx *svc.ServiceContext) func(context.Context, *DeletePaymentMethodInput) (*struct{}, error) {
	return func(ctx context.Context, input *DeletePaymentMethodInput) (*struct{}, error) {
		l := NewDeletePaymentMethodLogic(ctx, svcCtx)
		if err := l.DeletePaymentMethod(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type DeletePaymentMethodLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Delete Payment Method
func NewDeletePaymentMethodLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletePaymentMethodLogic {
	return &DeletePaymentMethodLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeletePaymentMethodLogic) DeletePaymentMethod(req *types.DeletePaymentMethodRequest) error {
	if err := l.svcCtx.PaymentModel.Delete(l.ctx, req.Id); err != nil {
		l.Errorw("delete payment method error", logger.Field("id", req.Id), logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseDeletedError), "delete payment method error: %s", err.Error())
	}
	return nil
}
