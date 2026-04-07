package smslogic

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/perfect-panel/server/modules/infra/logger"

	"github.com/hibiken/asynq"
	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/models/log"
	"github.com/perfect-panel/server/modules/notify/sms"
	"github.com/perfect-panel/server/worker/spec"
)

type SmsSendCount struct {
	Count    int   `json:"count"`
	CreateAt int64 `json:"create_at"`
}

type SendSmsLogic struct {
	deps Deps
}

func NewSendSmsLogic(deps Deps) *SendSmsLogic {
	return &SendSmsLogic{
		deps: deps,
	}
}
func (l *SendSmsLogic) ProcessTask(ctx context.Context, task *asynq.Task) error {
	var payload spec.SendSmsPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		logger.WithContext(ctx).Error("[SendSmsLogic] Unmarshal payload failed",
			logger.Field("error", err.Error()),
			logger.Field("payload", task.Payload()),
		)
		return nil
	}
	cfg := l.deps.currentConfig()
	client, err := sms.NewSender(cfg.Mobile.Platform, cfg.Mobile.PlatformConfig)
	if err != nil {
		logger.WithContext(ctx).Error("[SendSmsLogic] New send sms client failed", logger.Field("error", err.Error()), logger.Field("payload", payload))
		return err
	}
	createSms := &log.Message{
		Platform: cfg.Mobile.Platform,
		To:       fmt.Sprintf("+%s%s", payload.TelephoneArea, payload.Telephone),
		Subject:  config.ParseVerifyType(payload.Type).String(),
		Content: map[string]interface{}{
			"content": client.GetSendCodeContent(payload.Content),
		},
	}
	err = client.SendCode(payload.TelephoneArea, payload.Telephone, payload.Content)

	if err != nil {
		logger.WithContext(ctx).Error("[SendSmsLogic] Send sms failed", logger.Field("error", err.Error()), logger.Field("payload", payload))
		if cfg.Model != config.DevMode {
			createSms.Status = 2
		} else {
			return nil
		}
	}
	createSms.Status = 1
	logger.WithContext(ctx).Info("[SendSmsLogic] Send sms", logger.Field("telephone", payload.Telephone), logger.Field("content", createSms.Content))

	content, _ := createSms.Marshal()
	err = l.deps.LogModel.Insert(ctx, &log.SystemLog{
		Type:     log.TypeMobileMessage.Uint8(),
		Date:     time.Now().Format("2006-01-02"),
		ObjectID: 0,
		Content:  string(content),
	})
	if err != nil {
		logger.WithContext(ctx).Error("[SendSmsLogic] Send sms failed", logger.Field("error", err.Error()), logger.Field("payload", payload))
		return nil
	}
	return nil
}
