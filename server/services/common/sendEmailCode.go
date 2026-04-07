package common

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/modules/infra/limit"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/random"
	"github.com/perfect-panel/server/types"
	queue "github.com/perfect-panel/server/worker"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

type SendEmailCodeInput struct {
	Body types.SendCodeRequest
}

type SendEmailCodeOutput struct {
	Body *types.SendCodeResponse
}

func SendEmailCodeHandler(deps Deps) func(context.Context, *SendEmailCodeInput) (*SendEmailCodeOutput, error) {
	return func(ctx context.Context, input *SendEmailCodeInput) (*SendEmailCodeOutput, error) {
		l := NewSendEmailCodeLogic(ctx, deps)
		resp, err := l.SendEmailCode(&input.Body)
		if err != nil {
			return nil, err
		}
		return &SendEmailCodeOutput{Body: resp}, nil
	}
}

type SendEmailCodeLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

const (
	IntervalTime = 60
)

type VerifyTemplate struct {
	Type     uint8
	SiteLogo string
	SiteName string
	Expire   uint8
	Code     string
}
type CacheKeyPayload struct {
	Code   string `json:"code"`
	LastAt int64  `json:"lastAt"`
}

// NewSendEmailCodeLogic Get verification code
func NewSendEmailCodeLogic(ctx context.Context, deps Deps) *SendEmailCodeLogic {
	return &SendEmailCodeLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *SendEmailCodeLogic) SendEmailCode(req *types.SendCodeRequest) (resp *types.SendCodeResponse, err error) {
	cfg := l.deps.currentConfig()
	// Check if there is Redis in the code
	cacheKey := fmt.Sprintf("%s:%s:%s", config.AuthCodeCacheKey, config.ParseVerifyType(req.Type), req.Email)
	// Check if the limit is exceeded of current request
	limiter := limit.NewPeriodLimit(60, 1, l.deps.Redis, fmt.Sprintf("%s:%s:%s", config.SendIntervalKeyPrefix, "email", config.ParseVerifyType(req.Type)))
	permit, err := limiter.Take(req.Email)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "Failed to take limit")
	}
	if !limiter.ParsePermitState(permit) {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.TooManyRequests), "send email too many requests")
	}
	// Check if the limit is exceeded of today
	permit, err = l.deps.AuthLimiter.Take(fmt.Sprintf("%s:%s:%s", "email", config.ParseVerifyType(req.Type), req.Email))
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "Failed to take limit")
	}
	if !l.deps.AuthLimiter.ParsePermitState(permit) {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.TodaySendCountExceedsLimit), "send email too many requests")
	}
	m, err := l.deps.UserModel.FindUserAuthMethodByOpenID(l.ctx, "email", req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "FindUserAuthMethodByOpenID error")
	}
	if config.ParseVerifyType(req.Type) == config.Register && m.Id > 0 {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.UserExist), "mobile already bind")
	} else if config.ParseVerifyType(req.Type) == config.Security && m.Id == 0 {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.UserNotExist), "mobile not bind")
	}

	var payload CacheKeyPayload
	var taskPayload queue.SendEmailPayload
	// Generate verification code
	code := random.Key(6, 0)
	taskPayload.Type = queue.EmailTypeVerify
	taskPayload.Email = req.Email
	taskPayload.Subject = "Verification code"
	taskPayload.Content = map[string]interface{}{
		"Type":     req.Type,
		"SiteLogo": cfg.Site.SiteLogo,
		"SiteName": cfg.Site.SiteName,
		"Expire":   5,
		"Code":     code,
	}
	// Save to Redis
	payload = CacheKeyPayload{
		Code:   code,
		LastAt: time.Now().Unix(),
	}
	// Marshal the payload
	val, _ := json.Marshal(payload)
	if err = l.deps.Redis.Set(l.ctx, cacheKey, string(val), time.Second*IntervalTime*5).Err(); err != nil {
		l.Errorw("[SendEmailCode]: Redis Error", logger.Field("error", err.Error()), logger.Field("cacheKey", cacheKey))
		return nil, errors.Wrap(xerr.NewErrCode(xerr.ERROR), "Failed to set verification code")
	}

	// Marshal the task payload
	payloadBuy, err := json.Marshal(taskPayload)
	if err != nil {
		l.Errorw("[SendEmailCode]: Marshal Error", logger.Field("error", err.Error()))
		return nil, errors.Wrap(xerr.NewErrCode(xerr.ERROR), "Failed to marshal task payload")
	}
	// Create a queue task
	task := asynq.NewTask(queue.ForthwithSendEmail, payloadBuy, asynq.MaxRetry(3))
	// Enqueue the task
	taskInfo, err := l.deps.Queue.Enqueue(task)
	if err != nil {
		l.Errorw("[SendEmailCode]: Enqueue Error", logger.Field("error", err.Error()), logger.Field("payload", string(payloadBuy)))
		return nil, errors.Wrap(xerr.NewErrCode(xerr.ERROR), "Failed to enqueue task")
	}
	l.Infow("[SendEmailCode]: Enqueue Success", logger.Field("taskID", taskInfo.ID), logger.Field("payload", string(payloadBuy)))
	if cfg.Model == config.DevMode {
		return &types.SendCodeResponse{
			Code:   payload.Code,
			Status: true,
		}, nil
	} else {
		return &types.SendCodeResponse{
			Status: true,
		}, nil
	}
}
