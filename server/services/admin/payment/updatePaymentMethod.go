package payment

import (
	"context"
	"encoding/json"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/payment"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type UpdatePaymentMethodInput struct {
	Body types.UpdatePaymentMethodRequest
}

type UpdatePaymentMethodOutput struct {
	Body *types.PaymentConfig
}

func UpdatePaymentMethodHandler(deps Deps) func(context.Context, *UpdatePaymentMethodInput) (*UpdatePaymentMethodOutput, error) {
	return func(ctx context.Context, input *UpdatePaymentMethodInput) (*UpdatePaymentMethodOutput, error) {
		l := NewUpdatePaymentMethodLogic(ctx, deps)
		resp, err := l.UpdatePaymentMethod(&input.Body)
		if err != nil {
			return nil, err
		}
		return &UpdatePaymentMethodOutput{Body: resp}, nil
	}
}

type UpdatePaymentMethodLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// NewUpdatePaymentMethodLogic Update Payment Method
func NewUpdatePaymentMethodLogic(ctx context.Context, deps Deps) *UpdatePaymentMethodLogic {
	return &UpdatePaymentMethodLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *UpdatePaymentMethodLogic) UpdatePaymentMethod(req *types.UpdatePaymentMethodRequest) (resp *types.PaymentConfig, err error) {
	if payment.ParsePlatform(req.Platform) == payment.UNSUPPORTED {
		l.Errorw("unsupported payment platform", logger.Field("mark", req.Platform))
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(400, "UNSUPPORTED_PAYMENT_PLATFORM"), "unsupported payment platform: %s", req.Platform)
	}
	method, err := l.deps.PaymentModel.FindOne(l.ctx, req.Id)
	if err != nil {
		l.Errorw("find payment method error", logger.Field("id", req.Id), logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "find payment method error: %s", err.Error())
	}
	config := parsePaymentPlatformConfig(l.ctx, payment.ParsePlatform(req.Platform), req.Config)
	tool.DeepCopy(method, req, tool.CopyWithIgnoreEmpty(false))
	method.Config = config
	if err := l.deps.PaymentModel.Update(l.ctx, method); err != nil {
		l.Errorw("update payment method error", logger.Field("id", req.Id), logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseUpdateError), "update payment method error: %s", err.Error())
	}
	resp = &types.PaymentConfig{}
	tool.DeepCopy(resp, method)
	var configMap map[string]interface{}
	_ = json.Unmarshal([]byte(method.Config), &configMap)
	resp.Config = configMap
	return
}
