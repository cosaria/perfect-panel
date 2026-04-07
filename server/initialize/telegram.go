package initialize

import (
	"context"
	"fmt"

	"github.com/perfect-panel/server/modules/infra/logger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/models/auth"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/services/telegram"
)

func Telegram(deps Deps) {

	method, err := deps.AuthModel.FindOneByMethod(context.Background(), "telegram")
	if err != nil {
		logger.Errorf("[Init Telegram Config] Get Telegram Config Error: %s", err.Error())
		return
	}
	tgConfig := new(auth.TelegramAuthConfig)
	if err = tgConfig.Unmarshal(method.Config); err != nil {
		logger.Errorf("[Init Telegram Config] Unmarshal Telegram Config Error: %s", err.Error())
		return
	}
	tg := runtimeTelegramConfig(tgConfig)

	if tgConfig.BotToken == "" {
		clearTelegramRuntime(deps)
		logger.Debug("[Init Telegram Config] Telegram Token is empty")
		return
	}

	bot, err := tgbotapi.NewBotAPI(tg.BotToken)
	if err != nil {
		clearTelegramRuntime(deps)
		logger.Error("[Init Telegram Config] New Bot API Error: ", logger.Field("error", err.Error()))
		return
	}

	cfg := deps.currentConfig()
	if tgConfig.WebHookDomain == "" || cfg.Debug {
		// set Long Polling mode
		updateConfig := tgbotapi.NewUpdate(0)
		updateConfig.Timeout = 60
		updates := bot.GetUpdatesChan(updateConfig)
		replaceTelegramPoller(deps, bot.StopReceivingUpdates)
		go func() {
			deps := telegram.Deps{
				AuthModel:   deps.AuthModel,
				SystemModel: deps.SystemModel,
				UserModel:   deps.UserModel,
				Redis:       deps.Redis,
				DB:          deps.DB,
				TelegramBot: func() *tgbotapi.BotAPI { return bot },
				Config:      deps.Config,
			}
			for update := range updates {
				if update.Message != nil {
					ctx := context.Background()
					l := telegram.NewTelegramLogic(ctx, deps)
					l.TelegramLogic(&update)
				}
			}
		}()
	} else {
		replaceTelegramPoller(deps, nil)
		wh, err := tgbotapi.NewWebhook(fmt.Sprintf("%s/api/v1/telegram/webhook?secret=%s", tgConfig.WebHookDomain, tool.Md5Encode(tgConfig.BotToken, false)))
		if err != nil {
			clearTelegramRuntime(deps)
			logger.Errorf("[Init Telegram Config] New Webhook Error: %s", err.Error())
			return
		}
		_, err = bot.Request(wh)
		if err != nil {
			clearTelegramRuntime(deps)
			logger.Errorf("[Init Telegram Config] Request Webhook Error: %s", err.Error())
			return
		}
	}

	user, err := bot.GetMe()
	if err != nil {
		clearTelegramRuntime(deps)
		logger.Error("[Init Telegram Config] Get Bot Info Error: ", logger.Field("error", err.Error()))
		return
	}
	tg.BotID = user.ID
	tg.BotName = user.UserName
	syncTelegramRuntime(deps, tg, bot)

	logger.Info("[Init Telegram Config] Webhook set success")
}

func runtimeTelegramConfig(cfg *auth.TelegramAuthConfig) config.Telegram {
	if cfg == nil {
		return config.Telegram{}
	}
	return config.Telegram{
		Enable:        cfg.BotToken != "",
		BotToken:      cfg.BotToken,
		EnableNotify:  cfg.EnableNotify,
		WebHookDomain: cfg.WebHookDomain,
	}
}

func syncTelegramRuntime(deps Deps, tg config.Telegram, bot *tgbotapi.BotAPI) {
	if deps.Config != nil {
		deps.Config.Telegram = tg
	}
	if deps.SetTelegramBot != nil {
		deps.SetTelegramBot(bot)
	}
}

func replaceTelegramPoller(deps Deps, stop func()) {
	if deps.SwapTelegramPoller == nil {
		return
	}
	if previous := deps.SwapTelegramPoller(stop); previous != nil {
		previous()
	}
}

func clearTelegramRuntime(deps Deps) {
	replaceTelegramPoller(deps, nil)
	syncTelegramRuntime(deps, config.Telegram{}, nil)
}
