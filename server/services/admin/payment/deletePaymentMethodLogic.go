package payment

import (
	"context"

	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

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
