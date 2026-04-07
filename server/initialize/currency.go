package initialize

import (
	"context"
	"fmt"

	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/util/tool"
)

func Currency(deps Deps) {
	// Retrieve system currency configuration
	currency, err := deps.SystemModel.GetCurrencyConfig(context.Background())
	if err != nil {
		logger.Errorf("[INIT] Failed to get currency configuration: %v", err.Error())
		panic(fmt.Sprintf("[INIT] Failed to get currency configuration: %v", err.Error()))
	}
	// Parse currency configuration
	configs := struct {
		CurrencyUnit   string
		CurrencySymbol string
		AccessKey      string
	}{}
	tool.SystemConfigSliceReflectToStruct(currency, &configs)
	cfg := deps.currentConfig()
	cfg.Currency = config.Currency{
		Unit:      configs.CurrencyUnit,
		Symbol:    configs.CurrencySymbol,
		AccessKey: configs.AccessKey,
	}
	if deps.Config != nil {
		deps.Config.Currency = cfg.Currency
	}
	if deps.PrepareExchangeRate != nil {
		deps.PrepareExchangeRate(configs.CurrencyUnit, "CNY")
	} else if deps.SetExchangeRate != nil {
		deps.SetExchangeRate(0)
	}
	logger.Infof("[INIT] Currency configuration: %v", cfg.Currency)
}
