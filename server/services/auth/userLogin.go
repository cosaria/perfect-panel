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
	"github.com/perfect-panel/server/modules/verify/turnstile"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

type UserLoginInput struct {
	Body      types.UserLoginRequest
	IP        string `header:"X-Original-Forwarded-For" required:"false" doc:"Client IP from proxy"`
	UserAgent string `header:"User-Agent" required:"false" doc:"User agent string"`
	LoginType string `header:"Login-Type" required:"false" doc:"Login type"`
}

type UserLoginOutput struct {
	Body *types.LoginResponse
}

func UserLoginHandler(svcCtx *svc.ServiceContext) func(context.Context, *UserLoginInput) (*UserLoginOutput, error) {
	return func(ctx context.Context, input *UserLoginInput) (*UserLoginOutput, error) {
		input.Body.IP = input.IP
		input.Body.UserAgent = input.UserAgent
		input.Body.LoginType = input.LoginType
		if svcCtx.Config.Verify.LoginVerify && !svcCtx.Config.Debug {
			verifyTurns := turnstile.New(turnstile.Config{
				Secret:  svcCtx.Config.Verify.TurnstileSecret,
				Timeout: 3 * time.Second,
			})
			if verify, err := verifyTurns.Verify(ctx, input.Body.CfToken, input.Body.IP); err != nil || !verify {
				return nil, errors.Wrapf(xerr.NewErrCode(xerr.TooManyRequests), "error: %v, verify: %v", err, verify)
			}
		}
		l := NewUserLoginLogic(ctx, svcCtx)
		resp, err := l.UserLogin(&input.Body)
		if err != nil {
			return nil, err
		}
		return &UserLoginOutput{Body: resp}, nil
	}
}

type UserLoginLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewUserLoginLogic User login
func NewUserLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserLoginLogic {
	return &UserLoginLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserLoginLogic) UserLogin(req *types.UserLoginRequest) (resp *types.LoginResponse, err error) {
	loginStatus := false
	var userInfo *user.User
	// Record login status
	defer func(svcCtx *svc.ServiceContext) {
		if userInfo.Id != 0 {
			loginLog := log.Login{
				Method:    "email",
				LoginIP:   req.IP,
				UserAgent: req.UserAgent,
				Success:   loginStatus,
				Timestamp: time.Now().UnixMilli(),
			}
			content, _ := loginLog.Marshal()
			if err := l.svcCtx.LogModel.Insert(l.ctx, &log.SystemLog{
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
	}(l.svcCtx)

	userInfo, err = l.svcCtx.UserModel.FindOneByEmail(l.ctx, req.Email)

	if userInfo.DeletedAt.Valid {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.UserNotExist), "user email deleted: %v", req.Email)
	}

	if err != nil {
		if errors.As(err, &gorm.ErrRecordNotFound) {
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.UserNotExist), "user email not exist: %v", req.Email)
		}
		logger.WithContext(l.ctx).Error(err)
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "query user info failed: %v", err.Error())
	}

	// Verify password
	if !tool.MultiPasswordVerify(userInfo.Algo, userInfo.Salt, req.Password, userInfo.Password) {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.UserPasswordError), "user password")
	}

	// Bind device to user if identifier is provided
	if req.Identifier != "" {
		bindLogic := NewBindDeviceLogic(l.ctx, l.svcCtx)
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
		l.svcCtx.Config.JwtAuth.AccessSecret,
		time.Now().Unix(),
		l.svcCtx.Config.JwtAuth.AccessExpire,
		jwt.WithOption("UserId", userInfo.Id),
		jwt.WithOption("SessionId", sessionId),
		jwt.WithOption("LoginType", req.LoginType),
	)
	if err != nil {
		l.Logger.Error("[UserLogin] token generate error", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "token generate error: %v", err.Error())
	}
	sessionIdCacheKey := fmt.Sprintf("%v:%v", config.SessionIdKey, sessionId)
	if err = l.svcCtx.Redis.Set(l.ctx, sessionIdCacheKey, userInfo.Id, time.Duration(l.svcCtx.Config.JwtAuth.AccessExpire)*time.Second).Err(); err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "set session id error: %v", err.Error())
	}
	loginStatus = true
	return &types.LoginResponse{
		Token: token,
	}, nil
}
