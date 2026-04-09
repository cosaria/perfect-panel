package auth

import (
	"context"
	"encoding/json"
	"fmt"

	"time"

	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/models/log"
	"github.com/perfect-panel/server/models/user"
	"github.com/perfect-panel/server/modules/auth/jwt"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/modules/util/uuidx"
	"github.com/perfect-panel/server/modules/verify/turnstile"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type ResetPasswordInput struct {
	Body      types.ResetPasswordRequest
	IP        string `header:"X-Original-Forwarded-For" required:"false" doc:"Client IP from proxy"`
	UserAgent string `header:"User-Agent" required:"false" doc:"User agent string"`
	LoginType string `header:"Login-Type" required:"false" doc:"Login type"`
}

type ResetPasswordOutput struct {
	Body *types.LoginResponse
}

func ResetPasswordHandler(deps Deps) func(context.Context, *ResetPasswordInput) (*ResetPasswordOutput, error) {
	return func(ctx context.Context, input *ResetPasswordInput) (*ResetPasswordOutput, error) {
		input.Body.IP = input.IP
		input.Body.UserAgent = input.UserAgent
		input.Body.LoginType = input.LoginType
		cfg := deps.currentConfig()
		if cfg.Verify.ResetPasswordVerify {
			verifyTurns := turnstile.New(turnstile.Config{
				Secret:  cfg.Verify.TurnstileSecret,
				Timeout: 3 * time.Second,
			})
			if verify, err := verifyTurns.Verify(ctx, input.Body.CfToken, input.Body.IP); err != nil || !verify {
				return nil, errors.Wrapf(xerr.NewErrCode(xerr.TooManyRequests), "error: %v, verify: %v", err, verify)
			}
		}
		l := NewResetPasswordLogic(ctx, deps)
		resp, err := l.ResetPassword(&input.Body)
		if err != nil {
			return nil, err
		}
		return &ResetPasswordOutput{Body: resp}, nil
	}
}

type ResetPasswordLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// NewResetPasswordLogic Reset password
func NewResetPasswordLogic(ctx context.Context, deps Deps) *ResetPasswordLogic {
	return &ResetPasswordLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *ResetPasswordLogic) ResetPassword(req *types.ResetPasswordRequest) (resp *types.LoginResponse, err error) {
	var userInfo *user.User
	loginStatus := false

	defer func() {
		if userInfo.Id != 0 && loginStatus {
			loginLog := log.Login{
				Method:    "email",
				LoginIP:   req.IP,
				UserAgent: req.UserAgent,
				Success:   loginStatus,
				Timestamp: time.Now().UnixMilli(),
			}
			content, _ := loginLog.Marshal()
			if err := l.deps.LogModel.Insert(l.ctx, &log.SystemLog{
				Id:       0,
				Type:     log.TypeLogin.Uint8(),
				Date:     time.Now().Format("2006-01-02"),
				ObjectID: userInfo.Id,
				Content:  string(content),
			}); err != nil {
				l.Errorw("failed to insert login log",
					logger.Field("user_id", userInfo.Id),
					logger.Field("ip", req.IP),
					logger.Field("error", err.Error()),
				)
			}
		}
	}()

	cacheKey := fmt.Sprintf("%s:%s:%s", config.AuthCodeCacheKey, config.Security, req.Email)
	// Check the verification code
	if value, err := l.deps.Redis.Get(l.ctx, cacheKey).Result(); err != nil {
		l.Errorw("Verification code error", logger.Field("cacheKey", cacheKey), logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.VerifyCodeError), "Verification code error")
	} else {
		var payload CacheKeyPayload
		if err := json.Unmarshal([]byte(value), &payload); err != nil {
			l.Errorw("Unmarshal errors", logger.Field("cacheKey", cacheKey), logger.Field("error", err.Error()), logger.Field("value", value))
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.VerifyCodeError), "Verification code error")
		}
		if payload.Code != req.Code {
			l.Errorw("Verification code error", logger.Field("cacheKey", cacheKey), logger.Field("error", "Verification code error"), logger.Field("reqCode", req.Code), logger.Field("payloadCode", payload.Code))
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.VerifyCodeError), "Verification code error")
		}
	}

	// Check user
	authMethod, err := l.deps.UserModel.FindUserAuthMethodByOpenID(l.ctx, "email", req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.UserNotExist), "user email not exist: %v", req.Email)
		}
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "find user by email error: %v", err.Error())
	}

	userInfo, err = l.deps.UserModel.FindOne(l.ctx, authMethod.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.UserNotExist), "user email not exist: %v", req.Email)
		}
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "query user info failed: %v", err.Error())
	}

	// Update password
	userInfo.Password = tool.EncodePassWord(req.Password)
	userInfo.Algo = "default"
	if err = l.deps.UserModel.Update(l.ctx, userInfo); err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseUpdateError), "update user info failed: %v", err.Error())
	}

	// Bind device to user if identifier is provided
	if req.Identifier != "" {
		bindLogic := NewBindDeviceLogic(l.ctx, l.deps)
		if err := bindLogic.BindDeviceToUser(req.Identifier, req.IP, req.UserAgent, userInfo.Id); err != nil {
			l.Errorw("failed to bind device to user",
				logger.Field("user_id", userInfo.Id),
				logger.Field("identifier", req.Identifier),
				logger.Field("error", err.Error()),
			)
			// Don't fail register if device binding fails, just log the error
		}
	}
	if l.ctx.Value(config.LoginType) != nil {
		req.LoginType = l.ctx.Value(config.LoginType).(string)
	}
	// Generate session id
	sessionId := uuidx.NewUUID().String()
	// Generate token
	cfg := l.deps.currentConfig()
	token, err := jwt.NewJwtToken(
		cfg.JwtAuth.AccessSecret,
		time.Now().Unix(),
		cfg.JwtAuth.AccessExpire,
		jwt.WithOption("UserId", userInfo.Id),
		jwt.WithOption("SessionId", sessionId),
		jwt.WithOption("LoginType", req.LoginType),
	)
	if err != nil {
		l.Error("[UserLogin] token generate error", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "token generate error: %v", err.Error())
	}
	sessionIdCacheKey := fmt.Sprintf("%v:%v", config.SessionIdKey, sessionId)
	if err = l.deps.Redis.Set(l.ctx, sessionIdCacheKey, userInfo.Id, time.Duration(cfg.JwtAuth.AccessExpire)*time.Second).Err(); err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "set session id error: %v", err.Error())
	}
	loginStatus = true
	return &types.LoginResponse{
		Token: token,
	}, nil
}
