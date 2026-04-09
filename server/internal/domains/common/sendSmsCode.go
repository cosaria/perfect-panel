package common

import (
	"context"
	"encoding/json"
	"fmt"

	"time"

	"github.com/hibiken/asynq"
	"github.com/perfect-panel/server/config"
	queue "github.com/perfect-panel/server/internal/jobs"
	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/internal/platform/notify/phone"
	"github.com/perfect-panel/server/internal/platform/support/limit"
	"github.com/perfect-panel/server/internal/platform/support/logger"
	"github.com/perfect-panel/server/internal/platform/support/random"
	"github.com/perfect-panel/server/internal/platform/support/xerr"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type SendSmsCodeInput struct {
	Body types.SendSmsCodeRequest
}

type SendSmsCodeOutput struct {
	Body *types.SendCodeResponse
}

func SendSmsCodeHandler(deps Deps) func(context.Context, *SendSmsCodeInput) (*SendSmsCodeOutput, error) {
	return func(ctx context.Context, input *SendSmsCodeInput) (*SendSmsCodeOutput, error) {
		l := NewSendSmsCodeLogic(ctx, deps)
		resp, err := l.SendSmsCode(&input.Body)
		if err != nil {
			return nil, err
		}
		return &SendSmsCodeOutput{Body: resp}, nil
	}
}

type SmsSendCount struct {
	Count    int64 `json:"count"`
	CreateAt int64 `json:"create_at"`
}

type SendSmsCodeLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// NewSendSmsCodeLogic Get sms verification code
func NewSendSmsCodeLogic(ctx context.Context, deps Deps) *SendSmsCodeLogic {
	return &SendSmsCodeLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *SendSmsCodeLogic) SendSmsCode(req *types.SendSmsCodeRequest) (resp *types.SendCodeResponse, err error) {
	cfg := l.deps.currentConfig()
	phoneNumber, err := phone.FormatToE164(req.TelephoneAreaCode, req.Telephone)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.TelephoneError), "Invalid phone number")
	}

	cacheKey := fmt.Sprintf("%s:%s:%s", config.AuthCodeTelephoneCacheKey, config.ParseVerifyType(req.Type), phoneNumber)
	// Check if the limit is exceeded of current request
	limiter := limit.NewPeriodLimit(60, 1, l.deps.Redis, fmt.Sprintf("%s:%s:%s", config.SendIntervalKeyPrefix, "mobile", config.ParseVerifyType(req.Type)))
	permit, err := limiter.Take(phoneNumber)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "Failed to take limit")
	}
	if !limiter.ParsePermitState(permit) {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.TooManyRequests), "send sms too many requests")
	}
	// Check if the limit is exceeded of the today
	permit, err = l.deps.AuthLimiter.Take(fmt.Sprintf("%s:%s:%s", "mobile", config.ParseVerifyType(req.Type), phoneNumber))
	if err != nil {
		return nil, err
	}
	if !l.deps.AuthLimiter.ParsePermitState(permit) {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.TodaySendCountExceedsLimit), "This account has reached the limit of sending times today")
	}
	m, err := l.deps.UserModel.FindUserAuthMethodByOpenID(l.ctx, "mobile", phoneNumber)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "FindUserAuthMethodByOpenID error")
	}
	if config.ParseVerifyType(req.Type) == config.Register && m.Id > 0 {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.UserExist), "mobile already bind")
	} else if config.ParseVerifyType(req.Type) == config.Security && m.Id == 0 {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.UserNotExist), "mobile not bind")
	}

	taskPayload := queue.SendSmsPayload{
		Type:          req.Type,
		Telephone:     req.Telephone,
		TelephoneArea: req.TelephoneAreaCode,
	}
	// Generate verification code
	code := random.Key(6, 0)
	taskPayload.Telephone = req.Telephone
	taskPayload.Content = code
	// Save to Redis
	payload := CacheKeyPayload{
		Code:   code,
		LastAt: time.Now().Unix(),
	}
	// Marshal the payload
	val, _ := json.Marshal(payload)
	if err = l.deps.Redis.Set(l.ctx, cacheKey, string(val), time.Second*time.Duration(cfg.VerifyCode.ExpireTime)).Err(); err != nil {
		l.Errorw("[SendSmsCode]: Redis Error", logger.Field("error", err.Error()), logger.Field("cacheKey", cacheKey))
		return nil, errors.Wrap(xerr.NewErrCode(xerr.ERROR), "Failed to set verification code")
	}

	// Marshal the task payload
	payloadValue, err := json.Marshal(taskPayload)
	if err != nil {
		l.Errorw("[SendSmsCode]: Marshal Error", logger.Field("error", err.Error()))
		return nil, errors.Wrap(xerr.NewErrCode(xerr.ERROR), "Failed to marshal task payload")
	}
	// Create a queue task
	task := asynq.NewTask(queue.ForthwithSendSms, payloadValue)
	// Enqueue the task
	taskInfo, err := l.deps.Queue.Enqueue(task)
	if err != nil {
		l.Errorw("[SendSmsCode]: Enqueue Error", logger.Field("error", err.Error()), logger.Field("payload", string(payloadValue)))
		return nil, errors.Wrap(xerr.NewErrCode(xerr.ERROR), "Failed to enqueue task")
	}
	l.Infow("[SendSmsCode]: Enqueue Success", logger.Field("taskID", taskInfo.ID), logger.Field("payload", string(payloadValue)))
	if cfg.Model == config.DevMode {
		return &types.SendCodeResponse{
			Code:   taskPayload.Content,
			Status: true,
		}, nil
	}
	return &types.SendCodeResponse{
		Status: true,
	}, nil
}
