package payment

import (
	"context"

	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/internal/platform/support/logger"
	"github.com/perfect-panel/server/internal/platform/support/xerr"
	"github.com/pkg/errors"
)

type DeletePaymentMethodInput struct {
	Body types.DeletePaymentMethodRequest
}

func DeletePaymentMethodHandler(deps Deps) func(context.Context, *DeletePaymentMethodInput) (*struct{}, error) {
	return func(ctx context.Context, input *DeletePaymentMethodInput) (*struct{}, error) {
		l := NewDeletePaymentMethodLogic(ctx, deps)
		if err := l.DeletePaymentMethod(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type DeletePaymentMethodLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Delete Payment Method
func NewDeletePaymentMethodLogic(ctx context.Context, deps Deps) *DeletePaymentMethodLogic {
	return &DeletePaymentMethodLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *DeletePaymentMethodLogic) DeletePaymentMethod(req *types.DeletePaymentMethodRequest) error {
	if err := l.deps.PaymentModel.Delete(l.ctx, req.Id); err != nil {
		l.Errorw("delete payment method error", logger.Field("id", req.Id), logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseDeletedError), "delete payment method error: %s", err.Error())
	}
	return nil
}
