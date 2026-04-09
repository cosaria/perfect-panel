package user

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/models/user"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/notify/phone"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type UpdateBindMobileInput struct {
	Body types.UpdateBindMobileRequest
}

func UpdateBindMobileHandler(deps Deps) func(context.Context, *UpdateBindMobileInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateBindMobileInput) (*struct{}, error) {
		l := NewUpdateBindMobileLogic(ctx, deps)
		if err := l.UpdateBindMobile(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type UpdateBindMobileLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Update Bind Mobile
func NewUpdateBindMobileLogic(ctx context.Context, deps Deps) *UpdateBindMobileLogic {
	return &UpdateBindMobileLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *UpdateBindMobileLogic) UpdateBindMobile(req *types.UpdateBindMobileRequest) error {
	u, ok := l.ctx.Value(config.CtxKeyUser).(*user.User)
	if !ok {
		logger.Error("current user is not found in context")
		return errors.Wrapf(xerr.NewErrCode(xerr.InvalidAccess), "Invalid Access")
	}
	// verify mobile
	phoneNumber, err := phone.FormatToE164(req.AreaCode, req.Mobile)
	if err != nil {
		return errors.Wrapf(xerr.NewErrCode(xerr.TelephoneError), "Invalid phone number")
	}
	cacheKey := fmt.Sprintf("%s:%s:%s", config.AuthCodeTelephoneCacheKey, config.Register, phoneNumber)
	code, err := l.deps.Redis.Get(l.ctx, cacheKey).Result()
	if err != nil {
		l.Errorw("Redis Error", logger.Field("error", err.Error()), logger.Field("cacheKey", cacheKey))
		return errors.Wrapf(xerr.NewErrCode(xerr.VerifyCodeError), "code error")
	}
	var payload CacheKeyPayload
	err = json.Unmarshal([]byte(code), &payload)
	if err != nil {
		l.Errorw("Redis Error", logger.Field("error", err.Error()), logger.Field("cacheKey", cacheKey))
		return errors.Wrapf(xerr.NewErrCode(xerr.VerifyCodeError), "code error")
	}
	if payload.Code != req.Code {
		return errors.Wrapf(xerr.NewErrCode(xerr.VerifyCodeError), "code error")
	}
	l.deps.Redis.Del(l.ctx, cacheKey)

	m, err := l.deps.UserModel.FindUserAuthMethodByOpenID(l.ctx, "mobile", req.Mobile)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "FindUserAuthMethodByOpenID error")
	}
	if m.Id > 0 {
		return errors.Wrapf(xerr.NewErrCode(xerr.UserExist), "mobile already bind")
	}

	method, err := l.deps.UserModel.FindUserAuthMethodByUserId(l.ctx, "mobile", u.Id)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "FindUserAuthMethodByOpenID error")
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		method = &user.AuthMethods{
			UserId:         u.Id,
			AuthType:       "mobile",
			AuthIdentifier: req.Mobile,
			Verified:       true,
		}
		if err := l.deps.UserModel.InsertUserAuthMethods(l.ctx, method); err != nil {
			return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseInsertError), "InsertUserAuthMethods error")
		}
	} else {
		method.Verified = true
		method.AuthIdentifier = req.Mobile
		if err := l.deps.UserModel.UpdateUserAuthMethods(l.ctx, method); err != nil {
			return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseUpdateError), "UpdateUserAuthMethods error")
		}
	}
	return nil
}
