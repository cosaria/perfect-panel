package user

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/models/user"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/services/telegram"
	"github.com/perfect-panel/server/svc"
	"github.com/pkg/errors"
	"strconv"
	"time"
)

func UnbindTelegramHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*struct{}, error) {
	return func(ctx context.Context, _ *struct{}) (*struct{}, error) {
		l := NewUnbindTelegramLogic(ctx, svcCtx)
		if err := l.UnbindTelegram(); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type UnbindTelegramLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Unbind Telegram
func NewUnbindTelegramLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnbindTelegramLogic {
	return &UnbindTelegramLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnbindTelegramLogic) UnbindTelegram() error {
	// Get User Info
	u, ok := l.ctx.Value(config.CtxKeyUser).(*user.User)

	if !ok {
		logger.Error("current user is not found in context")
		return errors.Wrapf(xerr.NewErrCode(xerr.InvalidAccess), "Invalid Access")
	}
	method, err := l.svcCtx.UserModel.FindUserAuthMethodByPlatform(l.ctx, u.Id, "telegram")
	if err != nil {
		l.Errorw("UnbindTelegramLogic FindUserAuthMethodByPlatform Error", logger.Field("id", u.Id), logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "Find User Auth Method By Platform Failed")
	}

	userTelegramChatId, err := strconv.ParseInt(method.AuthIdentifier, 10, 64)
	if err != nil {
		l.Errorw("UnbindTelegramLogic ParseInt Error", logger.Field("id", u.Id), logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "ParseInt Error")
	}

	if userTelegramChatId == 0 {
		return errors.Wrapf(xerr.NewErrCode(xerr.TelegramNotBound), "Unbind Telegram")
	}

	// Unbind Telegram
	err = l.svcCtx.UserModel.DeleteUserAuthMethods(l.ctx, u.Id, "telegram")
	if err != nil {
		l.Errorw("UnbindTelegramLogic DeleteUserAuthMethods Error", logger.Field("id", u.Id), logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseDeletedError), "Delete User Auth Methods Failed")
	}
	// Unbind Telegram Success send message with chatId
	text, err := tool.RenderTemplateToString(telegram.UnbindNotify, map[string]string{
		"Id":   strconv.FormatInt(u.Id, 10),
		"Time": time.Now().Format("2006-01-02 15:04:05"),
	})
	if err != nil {
		l.Errorw("UnbindTelegramLogic RenderTemplateToString Error", logger.Field("id", u.Id), logger.Field("error", err.Error()))
		return nil
	}
	msg := tgbotapi.NewMessage(userTelegramChatId, text)
	_, err = l.svcCtx.TelegramBot.Send(msg)
	if err != nil {
		l.Errorw("UnbindTelegramLogic Send Error", logger.Field("id", u.Id), logger.Field("error", err.Error()))
		return nil
	}
	return nil
}
