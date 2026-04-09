package telegram

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/perfect-panel/server/internal/platform/support/logger"

	"github.com/gin-gonic/gin"
	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/internal/platform/persistence/auth"
	"github.com/perfect-panel/server/internal/platform/persistence/user"
	"github.com/perfect-panel/server/internal/platform/support/tool"
	"github.com/perfect-panel/server/internal/platform/support/xerr"
	"github.com/pkg/errors"
)

func GetTelegramConfig(ctx context.Context, deps Deps) (*types.TelegramConfig, error) {

	data, err := deps.AuthModel.FindOneByMethod(ctx, "telegram")
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "get Telegram config failed: %v", err.Error())
	}
	var telegramConfig auth.TelegramAuthConfig
	err = json.Unmarshal([]byte(data.Config), &telegramConfig)
	if err != nil {
		logger.WithContext(ctx).Error("unmarshal telegram config failed", logger.Field("error", err.Error()))
		return nil, err
	}

	return &types.TelegramConfig{
		TelegramBotToken:      telegramConfig.BotToken,
		TelegramNotify:        *data.Enabled,
		TelegramWebHookDomain: telegramConfig.WebHookDomain,
	}, nil
}

func ApiLink(ctx *gin.Context, deps Deps, method string) string {
	cfg, _ := GetTelegramConfig(ctx, deps)
	return "https://api.telegram.org/bot" + cfg.TelegramBotToken + "/" + method
}

func SendUserMessage(ctx *gin.Context, deps Deps, u user.User, text string, parseMode string) {
	req, _ := http.NewRequest("GET", ApiLink(ctx, deps, "sendMessage"), nil)
	q := req.URL.Query()

	userTelegramChatId, ok := findTelegram(&u)
	if !ok {
		return
	}
	q.Add("chat_id", strconv.FormatInt(userTelegramChatId, 10))
	if parseMode == "markdown" {
		text = strings.ReplaceAll(text, "_", "\\_")
	}
	q.Add("text", text)
	q.Add("parse_mode", parseMode)
	req.URL.RawQuery = q.Encode()
	_, _ = http.DefaultClient.Do(req)

}

func SendAdminMessage(ctx *gin.Context, deps Deps, text string, parseMode string) {
	var adminTelegram []int64
	f := false
	adminTelegramJson, err := deps.Redis.Get(ctx, "adminTelegram").Result()
	if err == nil {
		err = json.Unmarshal([]byte(adminTelegramJson), &adminTelegram)
		if err == nil {
			f = true
		}
	}
	if !f {
		deps.DB.Model(&user.User{}).Where("is_admin = true").Pluck("telegram", &adminTelegram)
		val, _ := json.Marshal(adminTelegram)
		_ = deps.Redis.Set(ctx, "TelegramConfig", string(val), time.Duration(3600)*time.Second).Err()
	}
	req, _ := http.NewRequest("GET", ApiLink(ctx, deps, "sendMessage"), nil)
	q := req.URL.Query()
	if parseMode == "markdown" {
		text = strings.ReplaceAll(text, "_", "\\_")
	}
	q.Add("text", text)
	q.Add("parse_mode", parseMode)
	for _, telegram := range adminTelegram {
		q.Add("chat_id", strconv.FormatInt(telegram, 10))
		req.URL.RawQuery = q.Encode()
		_, _ = http.DefaultClient.Do(req)
	}
}

func SetWebhook(ctx *gin.Context, deps Deps) error {
	configs, _ := deps.SystemModel.GetSiteConfig(ctx)
	cfg := &types.SiteConfig{}
	tool.SystemConfigSliceReflectToStruct(configs, cfg)
	req, _ := http.NewRequest("GET", ApiLink(ctx, deps, "setWebhook"), nil)
	q := req.URL.Query()
	q.Add("url", cfg.Host+"/telegram/webhook")
	req.URL.RawQuery = q.Encode()
	_, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "set webhook error: %v", err)
	}
	return nil
}

func findTelegram(u *user.User) (int64, bool) {
	for _, item := range u.AuthMethods {
		if item.AuthType == "telegram" {
			// string to int64
			parseInt, err := strconv.ParseInt(item.AuthIdentifier, 10, 64)
			if err != nil {
				return 0, false
			}
			return parseInt, true
		}

	}
	return 0, false
}
