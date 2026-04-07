package task

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/payment/exchangeRate"
	"github.com/perfect-panel/server/modules/util/tool"
)

type RateLogic struct {
	deps Deps
}

func NewRateLogic(deps Deps) *RateLogic {
	return &RateLogic{
		deps: deps,
	}
}

func (l *RateLogic) ProcessTask(ctx context.Context, _ *asynq.Task) error {
	// Retrieve system currency configuration
	currency, err := l.deps.SystemModel.GetCurrencyConfig(ctx)
	if err != nil {
		logger.Errorw("[PurchaseCheckout] GetCurrencyConfig error", logger.Field("error", err.Error()))
		return err
	}
	// Parse currency configuration
	configs := struct {
		CurrencyUnit   string
		CurrencySymbol string
		AccessKey      string
	}{}
	tool.SystemConfigSliceReflectToStruct(currency, &configs)

	// Skip conversion if no exchange rate API key configured
	if configs.AccessKey == "" {
		logger.Debugf("[RateLogic] skip exchange rate, no access key configured")
		return nil
	}
	// Update exchange rates
	result, err := exchangeRate.GetExchangeRete(configs.CurrencyUnit, "CNY", configs.AccessKey, 1)
	if err != nil {
		logger.Errorw("[RateLogic] GetExchangeRete error", logger.Field("error", err.Error()))
		return err
	}
	if l.deps.SetExchangeRate != nil {
		l.deps.SetExchangeRate(result)
	}
	logger.WithContext(ctx).Infof("[RateLogic] GetExchangeRete success, result: %+v", result)
	return nil
}
