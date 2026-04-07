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
	var tg config.Telegram

	tgConfig := new(auth.TelegramAuthConfig)
	if err = tgConfig.Unmarshal(method.Config); err != nil {
		logger.Errorf("[Init Telegram Config] Unmarshal Telegram Config Error: %s", err.Error())
		return
	}

	if tgConfig.BotToken == "" {
		logger.Debug("[Init Telegram Config] Telegram Token is empty")
		return
	}

	bot, err := tgbotapi.NewBotAPI(tg.BotToken)
	if err != nil {
		logger.Error("[Init Telegram Config] New Bot API Error: ", logger.Field("error", err.Error()))
		return
	}

	cfg := deps.currentConfig()
	if tgConfig.WebHookDomain == "" || cfg.Debug {
		// set Long Polling mode
		updateConfig := tgbotapi.NewUpdate(0)
		updateConfig.Timeout = 60
		updates := bot.GetUpdatesChan(updateConfig)
		go func() {
			deps := telegram.Deps{
				AuthModel:   deps.AuthModel,
				SystemModel: deps.SystemModel,
				UserModel:   deps.UserModel,
				Redis:       deps.Redis,
				DB:          deps.DB,
				TelegramBot: bot,
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
		wh, err := tgbotapi.NewWebhook(fmt.Sprintf("%s/api/v1/telegram/webhook?secret=%s", tgConfig.WebHookDomain, tool.Md5Encode(tgConfig.BotToken, false)))
		if err != nil {
			logger.Errorf("[Init Telegram Config] New Webhook Error: %s", err.Error())
			return
		}
		_, err = bot.Request(wh)
		if err != nil {
			logger.Errorf("[Init Telegram Config] Request Webhook Error: %s", err.Error())
			return
		}
	}

	user, err := bot.GetMe()
	if err != nil {
		logger.Error("[Init Telegram Config] Get Bot Info Error: ", logger.Field("error", err.Error()))
		return
	}
	if deps.Config != nil {
		deps.Config.Telegram.BotID = user.ID
		deps.Config.Telegram.BotName = user.UserName
		deps.Config.Telegram.EnableNotify = tg.EnableNotify
		deps.Config.Telegram.WebHookDomain = tg.WebHookDomain
	}
	if deps.SetTelegramBot != nil {
		deps.SetTelegramBot(bot)
	}

	logger.Info("[Init Telegram Config] Webhook set success")
}
