package auth

import (
	"context"
	"encoding/json"
	"fmt"

	"time"

	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/internal/platform/persistence/log"
	"github.com/perfect-panel/server/internal/platform/persistence/user"
	"github.com/perfect-panel/server/internal/platform/support/auth/jwt"
	"github.com/perfect-panel/server/internal/platform/support/logger"
	"github.com/perfect-panel/server/internal/platform/support/tool"
	"github.com/perfect-panel/server/internal/platform/support/uuidx"
	"github.com/perfect-panel/server/internal/platform/support/verify/turnstile"
	"github.com/perfect-panel/server/internal/platform/support/xerr"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type UserRegisterInput struct {
	Body      types.UserRegisterRequest
	IP        string `header:"X-Original-Forwarded-For" required:"false" doc:"Client IP from proxy"`
	UserAgent string `header:"User-Agent" required:"false" doc:"User agent string"`
	LoginType string `header:"Login-Type" required:"false" doc:"Login type"`
}

type UserRegisterOutput struct {
	Body *types.LoginResponse
}

func UserRegisterHandler(deps Deps) func(context.Context, *UserRegisterInput) (*UserRegisterOutput, error) {
	return func(ctx context.Context, input *UserRegisterInput) (*UserRegisterOutput, error) {
		input.Body.IP = input.IP
		input.Body.UserAgent = input.UserAgent
		input.Body.LoginType = input.LoginType
		cfg := deps.currentConfig()
		if cfg.Verify.RegisterVerify {
			verifyTurns := turnstile.New(turnstile.Config{
				Secret:  cfg.Verify.TurnstileSecret,
				Timeout: 3 * time.Second,
			})
			if verify, err := verifyTurns.Verify(ctx, input.Body.CfToken, input.Body.IP); err != nil || !verify {
				return nil, errors.Wrapf(xerr.NewErrCode(xerr.TooManyRequests), "verify error: %v", err)
			}
		}
		l := NewUserRegisterLogic(ctx, deps)
		resp, err := l.UserRegister(&input.Body)
		if err != nil {
			return nil, err
		}
		return &UserRegisterOutput{Body: resp}, nil
	}
}

type UserRegisterLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// NewUserRegisterLogic User register
func NewUserRegisterLogic(ctx context.Context, deps Deps) *UserRegisterLogic {
	return &UserRegisterLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *UserRegisterLogic) UserRegister(req *types.UserRegisterRequest) (resp *types.LoginResponse, err error) {
	cfg := l.deps.currentConfig()
	c := cfg.Register
	email := cfg.Email
	var referer *user.User
	// Check if the registration is stopped
	if c.StopRegister {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.StopRegister), "stop register")
	}

	if req.Invite == "" {
		if cfg.Invite.ForcedInvite {
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.InviteCodeError), "invite code is required")
		}
	} else {
		// Check if the invite code is valid
		referer, err = l.deps.UserModel.FindOneByReferCode(l.ctx, req.Invite)
		if err != nil {
			l.Errorw("FindOneByReferCode Error", logger.Field("error", err))
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.InviteCodeError), "invite code is invalid")
		}
	}

	// if the email verification is enabled, the verification code is required
	if email.EnableVerify {
		cacheKey := fmt.Sprintf("%s:%s:%s", config.AuthCodeCacheKey, config.Register, req.Email)
		value, err := l.deps.Redis.Get(l.ctx, cacheKey).Result()
		if err != nil {
			l.Errorw("Redis Error", logger.Field("error", err.Error()), logger.Field("cacheKey", cacheKey))
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.VerifyCodeError), "code error")
		}
		var payload CacheKeyPayload
		err = json.Unmarshal([]byte(value), &payload)
		if err != nil {
			l.Errorw("Unmarshal Error", logger.Field("error", err.Error()), logger.Field("value", value))
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.VerifyCodeError), "code error")
		}
		if payload.Code != req.Code {
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.VerifyCodeError), "code error")
		}
	}
	// Check if the user exists
	u, err := l.deps.UserModel.FindOneByEmail(l.ctx, req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		l.Errorw("FindOneByEmail Error", logger.Field("error", err))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "query user info failed: %v", err.Error())
	} else if err == nil && !u.DeletedAt.Valid {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.UserExist), "user email exist: %v", req.Email)
	} else if err == nil && u.DeletedAt.Valid {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.UserDisabled), "user email deleted: %v", req.Email)
	}

	// Generate password
	pwd := tool.EncodePassWord(req.Password)
	userInfo := &user.User{
		Password:          pwd,
		Algo:              "default",
		OnlyFirstPurchase: &cfg.Invite.OnlyFirstPurchase,
	}
	if referer != nil {
		userInfo.RefererId = referer.Id
	}
	err = l.deps.UserModel.Transaction(l.ctx, func(db *gorm.DB) error {
		// Save user information
		if err := l.deps.UserModel.Insert(l.ctx, userInfo, db); err != nil {
			return err
		}
		// Generate ReferCode
		userInfo.ReferCode = uuidx.UserInviteCode(userInfo.Id)
		// Update ReferCode
		if err := l.deps.UserModel.Update(l.ctx, userInfo, db); err != nil {
			return err
		}
		// create user auth info
		authInfo := &user.AuthMethods{
			UserId:         userInfo.Id,
			AuthType:       "email",
			AuthIdentifier: req.Email,
			Verified:       email.EnableVerify,
		}
		if err = l.deps.UserModel.InsertUserAuthMethods(l.ctx, authInfo, db); err != nil {
			return err
		}

		if cfg.Register.EnableTrial {
			// Active trial
			if err = l.activeTrial(userInfo.Id); err != nil {
				return err
			}
		}
		return nil
	})
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
	// Set session id
	sessionIdCacheKey := fmt.Sprintf("%v:%v", config.SessionIdKey, sessionId)
	if err := l.deps.Redis.Set(l.ctx, sessionIdCacheKey, userInfo.Id, time.Duration(cfg.JwtAuth.AccessExpire)*time.Second).Err(); err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "set session id error: %v", err.Error())
	}
	loginStatus := true
	defer func() {
		if token != "" && userInfo.Id != 0 {
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

			// Register log
			registerLog := log.Register{
				AuthMethod: "email",
				Identifier: req.Email,
				RegisterIP: req.IP,
				UserAgent:  req.UserAgent,
				Timestamp:  time.Now().UnixMilli(),
			}
			content, _ = registerLog.Marshal()
			if err = l.deps.LogModel.Insert(l.ctx, &log.SystemLog{
				Type:     log.TypeRegister.Uint8(),
				ObjectID: userInfo.Id,
				Date:     time.Now().Format("2006-01-02"),
				Content:  string(content),
			}); err != nil {
				l.Errorw("failed to insert login log",
					logger.Field("user_id", userInfo.Id),
					logger.Field("ip", req.IP),
					logger.Field("error", err.Error()))
			}
		}
	}()
	return &types.LoginResponse{
		Token: token,
	}, nil
}

func (l *UserRegisterLogic) activeTrial(uid int64) error {
	cfg := l.deps.currentConfig()
	sub, err := l.deps.SubscribeModel.FindOne(l.ctx, cfg.Register.TrialSubscribe)
	if err != nil {
		return err
	}
	userSub := &user.Subscribe{
		UserId:      uid,
		OrderId:     0,
		SubscribeId: sub.Id,
		StartTime:   time.Now(),
		ExpireTime:  tool.AddTime(cfg.Register.TrialTimeUnit, cfg.Register.TrialTime, time.Now()),
		Traffic:     sub.Traffic,
		Download:    0,
		Upload:      0,
		Token:       uuidx.SubscribeToken(fmt.Sprintf("Trial-%v", uid)),
		UUID:        uuidx.NewUUID().String(),
		Status:      1,
	}
	return l.deps.UserModel.InsertSubscribe(l.ctx, userSub)
}
