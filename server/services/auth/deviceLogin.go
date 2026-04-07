package auth

import (
	"context"
	"fmt"

	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/models/log"
	"github.com/perfect-panel/server/models/user"
	"github.com/perfect-panel/server/modules/auth/jwt"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/modules/util/uuidx"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

type DeviceLoginInput struct {
	Body types.DeviceLoginRequest
	IP   string `header:"X-Original-Forwarded-For" required:"false" doc:"Client IP from proxy"`
}

type DeviceLoginOutput struct {
	Body *types.LoginResponse
}

func DeviceLoginHandler(deps Deps) func(context.Context, *DeviceLoginInput) (*DeviceLoginOutput, error) {
	return func(ctx context.Context, input *DeviceLoginInput) (*DeviceLoginOutput, error) {
		input.Body.IP = input.IP
		l := NewDeviceLoginLogic(ctx, deps)
		resp, err := l.DeviceLogin(&input.Body)
		if err != nil {
			return nil, err
		}
		return &DeviceLoginOutput{Body: resp}, nil
	}
}

type DeviceLoginLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Device Login
func NewDeviceLoginLogic(ctx context.Context, deps Deps) *DeviceLoginLogic {
	return &DeviceLoginLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *DeviceLoginLogic) DeviceLogin(req *types.DeviceLoginRequest) (resp *types.LoginResponse, err error) {
	cfg := l.deps.currentConfig()
	if !cfg.Device.Enable {
		return nil, xerr.NewErrMsg("Device login is disabled")
	}

	loginStatus := false
	var userInfo *user.User
	// Record login status
	defer func() {
		if userInfo != nil && userInfo.Id != 0 {
			loginLog := log.Login{
				Method:    "device",
				LoginIP:   req.IP,
				UserAgent: req.UserAgent,
				Success:   loginStatus,
				Timestamp: time.Now().UnixMilli(),
			}
			content, _ := loginLog.Marshal()
			if err := l.deps.LogModel.Insert(l.ctx, &log.SystemLog{
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

	// Check if device exists by identifier
	deviceInfo, err := l.deps.UserModel.FindOneDeviceByIdentifier(l.ctx, req.Identifier)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Device not found, create new user and device
			userInfo, err = l.registerUserAndDevice(req)
			if err != nil {
				return nil, err
			}
		} else {
			l.Errorw("query device failed",
				logger.Field("identifier", req.Identifier),
				logger.Field("error", err.Error()),
			)
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "query device failed: %v", err.Error())
		}
	} else {
		// Device found, get user info
		userInfo, err = l.deps.UserModel.FindOne(l.ctx, deviceInfo.UserId)
		if err != nil {
			l.Errorw("query user failed",
				logger.Field("user_id", deviceInfo.UserId),
				logger.Field("error", err.Error()),
			)
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "query user failed: %v", err.Error())
		}
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
		jwt.WithOption("LoginType", "device"),
	)
	if err != nil {
		l.Errorw("token generate error",
			logger.Field("user_id", userInfo.Id),
			logger.Field("error", err.Error()),
		)
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "token generate error: %v", err.Error())
	}

	// Store session id in redis
	sessionIdCacheKey := fmt.Sprintf("%v:%v", config.SessionIdKey, sessionId)
	if err = l.deps.Redis.Set(l.ctx, sessionIdCacheKey, userInfo.Id, time.Duration(cfg.JwtAuth.AccessExpire)*time.Second).Err(); err != nil {
		l.Errorw("set session id error",
			logger.Field("user_id", userInfo.Id),
			logger.Field("error", err.Error()),
		)
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "set session id error: %v", err.Error())
	}

	loginStatus = true
	return &types.LoginResponse{
		Token: token,
	}, nil
}

func (l *DeviceLoginLogic) registerUserAndDevice(req *types.DeviceLoginRequest) (*user.User, error) {
	l.Infow("device not found, creating new user and device",
		logger.Field("identifier", req.Identifier),
		logger.Field("ip", req.IP),
	)

	var userInfo *user.User
	cfg := l.deps.currentConfig()
	err := l.deps.UserModel.Transaction(l.ctx, func(db *gorm.DB) error {
		// Create new user
		userInfo = &user.User{
			OnlyFirstPurchase: &cfg.Invite.OnlyFirstPurchase,
		}
		if err := db.Create(userInfo).Error; err != nil {
			l.Errorw("failed to create user",
				logger.Field("error", err.Error()),
			)
			return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseInsertError), "create user failed: %v", err)
		}

		// Update refer code
		userInfo.ReferCode = uuidx.UserInviteCode(userInfo.Id)
		if err := db.Model(&user.User{}).Where("id = ?", userInfo.Id).Update("refer_code", userInfo.ReferCode).Error; err != nil {
			l.Errorw("failed to update refer code",
				logger.Field("user_id", userInfo.Id),
				logger.Field("error", err.Error()),
			)
			return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseUpdateError), "update refer code failed: %v", err)
		}

		// Create device auth method
		authMethod := &user.AuthMethods{
			UserId:         userInfo.Id,
			AuthType:       "device",
			AuthIdentifier: req.Identifier,
			Verified:       true,
		}
		if err := db.Create(authMethod).Error; err != nil {
			l.Errorw("failed to create device auth method",
				logger.Field("user_id", userInfo.Id),
				logger.Field("identifier", req.Identifier),
				logger.Field("error", err.Error()),
			)
			return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseInsertError), "create device auth method failed: %v", err)
		}

		// Insert device record
		deviceInfo := &user.Device{
			Ip:         req.IP,
			UserId:     userInfo.Id,
			UserAgent:  req.UserAgent,
			Identifier: req.Identifier,
			Enabled:    true,
			Online:     false,
		}
		if err := db.Create(deviceInfo).Error; err != nil {
			l.Errorw("failed to insert device",
				logger.Field("user_id", userInfo.Id),
				logger.Field("identifier", req.Identifier),
				logger.Field("error", err.Error()),
			)
			return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseInsertError), "insert device failed: %v", err)
		}

		// Activate trial if enabled
		if cfg.Register.EnableTrial {
			if err := l.activeTrial(userInfo.Id, db); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		l.Errorw("device registration failed",
			logger.Field("identifier", req.Identifier),
			logger.Field("error", err.Error()),
		)
		return nil, err
	}

	l.Infow("device registration completed successfully",
		logger.Field("user_id", userInfo.Id),
		logger.Field("identifier", req.Identifier),
		logger.Field("refer_code", userInfo.ReferCode),
	)

	// Register log
	registerLog := log.Register{
		AuthMethod: "device",
		Identifier: req.Identifier,
		RegisterIP: req.IP,
		UserAgent:  req.UserAgent,
		Timestamp:  time.Now().UnixMilli(),
	}
	content, _ := registerLog.Marshal()

	if err := l.deps.LogModel.Insert(l.ctx, &log.SystemLog{
		Type:     log.TypeRegister.Uint8(),
		Date:     time.Now().Format("2006-01-02"),
		ObjectID: userInfo.Id,
		Content:  string(content),
	}); err != nil {
		l.Errorw("failed to insert register log",
			logger.Field("user_id", userInfo.Id),
			logger.Field("ip", req.IP),
			logger.Field("error", err.Error()),
		)
	}

	return userInfo, nil
}

func (l *DeviceLoginLogic) activeTrial(userId int64, db *gorm.DB) error {
	cfg := l.deps.currentConfig()
	sub, err := l.deps.SubscribeModel.FindOne(l.ctx, cfg.Register.TrialSubscribe)
	if err != nil {
		l.Errorw("failed to find trial subscription template",
			logger.Field("user_id", userId),
			logger.Field("trial_subscribe_id", cfg.Register.TrialSubscribe),
			logger.Field("error", err.Error()),
		)
		return err
	}

	startTime := time.Now()
	expireTime := tool.AddTime(cfg.Register.TrialTimeUnit, cfg.Register.TrialTime, startTime)
	subscribeToken := uuidx.SubscribeToken(fmt.Sprintf("Trial-%v", userId))
	subscribeUUID := uuidx.NewUUID().String()

	userSub := &user.Subscribe{
		UserId:      userId,
		OrderId:     0,
		SubscribeId: sub.Id,
		StartTime:   startTime,
		ExpireTime:  expireTime,
		Traffic:     sub.Traffic,
		Download:    0,
		Upload:      0,
		Token:       subscribeToken,
		UUID:        subscribeUUID,
		Status:      1,
	}

	if err := db.Create(userSub).Error; err != nil {
		l.Errorw("failed to insert trial subscription",
			logger.Field("user_id", userId),
			logger.Field("error", err.Error()),
		)
		return err
	}

	l.Infow("trial subscription activated successfully",
		logger.Field("user_id", userId),
		logger.Field("subscribe_id", sub.Id),
		logger.Field("expire_time", expireTime),
		logger.Field("traffic", sub.Traffic),
	)

	return nil
}
