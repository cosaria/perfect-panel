package auth

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/models/log"
	"github.com/perfect-panel/server/models/user"
	"github.com/perfect-panel/server/modules/auth/jwt"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/notify/phone"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/modules/util/uuidx"
	"github.com/perfect-panel/server/modules/verify/turnstile"
	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type TelephoneLoginInput struct {
	Body      types.TelephoneLoginRequest
	IP        string `header:"X-Original-Forwarded-For" required:"false" doc:"Client IP from proxy"`
	UserAgent string `header:"User-Agent" required:"false" doc:"User agent string"`
	LoginType string `header:"Login-Type" required:"false" doc:"Login type"`
}

type TelephoneLoginOutput struct {
	Body *types.LoginResponse
}

func TelephoneLoginHandler(deps Deps) func(context.Context, *TelephoneLoginInput) (*TelephoneLoginOutput, error) {
	return func(ctx context.Context, input *TelephoneLoginInput) (*TelephoneLoginOutput, error) {
		input.Body.IP = input.IP
		input.Body.UserAgent = input.UserAgent
		input.Body.LoginType = input.LoginType
		cfg := deps.currentConfig()
		if cfg.Verify.LoginVerify {
			verifyTurns := turnstile.New(turnstile.Config{
				Secret:  cfg.Verify.TurnstileSecret,
				Timeout: 3 * time.Second,
			})
			if verify, err := verifyTurns.Verify(ctx, input.Body.CfToken, input.Body.IP); err != nil || !verify {
				return nil, errors.Wrapf(xerr.NewErrCode(xerr.TooManyRequests), "error: %v, verify: %v", err, verify)
			}
		}
		// Construct a minimal *http.Request with User-Agent header for the logic layer
		r := &http.Request{Header: http.Header{}}
		r.Header.Set("User-Agent", input.UserAgent)
		l := NewTelephoneLoginLogic(ctx, deps)
		resp, err := l.TelephoneLogin(&input.Body, r, input.IP)
		if err != nil {
			return nil, err
		}
		return &TelephoneLoginOutput{Body: resp}, nil
	}
}

type TelephoneLoginLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// User Telephone login
func NewTelephoneLoginLogic(ctx context.Context, deps Deps) *TelephoneLoginLogic {
	return &TelephoneLoginLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *TelephoneLoginLogic) TelephoneLogin(req *types.TelephoneLoginRequest, r *http.Request, ip string) (resp *types.LoginResponse, err error) {
	phoneNumber, err := phone.FormatToE164(req.TelephoneAreaCode, req.Telephone)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.TelephoneError), "Invalid phone number")
	}
	cfg := l.deps.currentConfig()
	if !cfg.Mobile.Enable {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.SmsNotEnabled), "sms login is not enabled")
	}
	loginStatus := false
	var userInfo *user.User
	// Record login status
	defer func() {
		if userInfo.Id != 0 {
			loginLog := log.Login{
				Method:    "mobile",
				LoginIP:   ip,
				UserAgent: r.UserAgent(),
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

	authMethodInfo, err := l.deps.UserModel.FindUserAuthMethodByOpenID(l.ctx, "mobile", phoneNumber)
	if err != nil {
		if errors.As(err, gorm.ErrRecordNotFound) {
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.UserNotExist), "user telephone not exist: %v", req.Telephone)
		}
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "query user info failed: %v", err.Error())
	}

	userInfo, err = l.deps.UserModel.FindOne(l.ctx, authMethodInfo.UserId)
	if err != nil {
		if errors.As(err, gorm.ErrRecordNotFound) {
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.UserNotExist), "user telephone not exist: %v", req.Telephone)
		}
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "query user info failed: %v", err.Error())
	}

	if req.Password == "" && req.TelephoneCode == "" {
		return nil, xerr.NewErrCodeMsg(xerr.InvalidParams, "password and telephone code is empty")
	}

	if req.TelephoneCode == "" {
		// Verify password
		if !tool.MultiPasswordVerify(userInfo.Algo, userInfo.Salt, req.Password, userInfo.Password) {
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.UserPasswordError), "user password")
		}
	} else {
		cacheKey := fmt.Sprintf("%s:%s:%s", config.AuthCodeTelephoneCacheKey, config.ParseVerifyType(uint8(config.Security)), phoneNumber)
		value, err := l.deps.Redis.Get(l.ctx, cacheKey).Result()
		if err != nil {
			l.Errorw("Redis Error", logger.Field("error", err.Error()), logger.Field("cacheKey", cacheKey))
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.VerifyCodeError), "code error")
		}

		if value == "" {
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.VerifyCodeError), "code error")
		}

		var payload CacheKeyPayload
		if err := json.Unmarshal([]byte(value), &payload); err != nil {
			l.Errorw("[SendSmsCode]: Unmarshal Error", logger.Field("error", err.Error()), logger.Field("value", value))
		}

		if payload.Code != req.TelephoneCode {
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.VerifyCodeError), "code error")
		}
		l.deps.Redis.Del(l.ctx, cacheKey)
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
			// Don't fail login if device binding fails, just log the error
		}
	}

	if l.ctx.Value(config.LoginType) != nil {
		req.LoginType = l.ctx.Value(config.LoginType).(string)
	}

	// Generate session id
	sessionId := uuidx.NewUUID().String()
	// Generate token
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
