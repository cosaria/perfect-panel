package user

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/models/user"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type VerifyEmailInput struct {
	Body types.VerifyEmailRequest
}

func VerifyEmailHandler(svcCtx *svc.ServiceContext) func(context.Context, *VerifyEmailInput) (*struct{}, error) {
	return func(ctx context.Context, input *VerifyEmailInput) (*struct{}, error) {
		l := NewVerifyEmailLogic(ctx, svcCtx)
		if err := l.VerifyEmail(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type VerifyEmailLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Verify Email
func NewVerifyEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyEmailLogic {
	return &VerifyEmailLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

type CacheKeyPayload struct {
	Code   string `json:"code"`
	LastAt int64  `json:"lastAt"`
}

func (l *VerifyEmailLogic) VerifyEmail(req *types.VerifyEmailRequest) error {
	cacheKey := fmt.Sprintf("%s:%s:%s", config.AuthCodeCacheKey, config.Security, req.Email)
	value, err := l.svcCtx.Redis.Get(l.ctx, cacheKey).Result()
	if err != nil {
		l.Errorw("Redis Error", logger.Field("error", err.Error()), logger.Field("cacheKey", cacheKey))
		return errors.Wrapf(xerr.NewErrCode(xerr.VerifyCodeError), "code error")
	}

	var payload CacheKeyPayload
	err = json.Unmarshal([]byte(value), &payload)
	if err != nil {
		l.Errorw("Redis Error", logger.Field("error", err.Error()), logger.Field("cacheKey", cacheKey))
		return errors.Wrapf(xerr.NewErrCode(xerr.VerifyCodeError), "code error")
	}
	if payload.Code != req.Code {
		return errors.Wrapf(xerr.NewErrCode(xerr.VerifyCodeError), "code error")
	}
	l.svcCtx.Redis.Del(l.ctx, cacheKey)

	u, ok := l.ctx.Value(config.CtxKeyUser).(*user.User)
	if !ok {
		logger.Error("current user is not found in context")
		return errors.Wrapf(xerr.NewErrCode(xerr.InvalidAccess), "Invalid Access")
	}
	method, err := l.svcCtx.UserModel.FindUserAuthMethodByOpenID(l.ctx, "email", req.Email)
	if err != nil {
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "FindUserAuthMethodByOpenID error")
	}
	if method.UserId != u.Id {
		return errors.Wrapf(xerr.NewErrCode(xerr.InvalidAccess), "invalid access")
	}
	method.Verified = true
	err = l.svcCtx.UserModel.UpdateUserAuthMethods(l.ctx, method)
	if err != nil {
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseUpdateError), "UpdateUserAuthMethods error")
	}
	return nil
}
