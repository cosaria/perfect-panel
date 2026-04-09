package common

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/internal/platform/notify/phone"
	"github.com/perfect-panel/server/internal/platform/support/auth/authmethod"
	"github.com/perfect-panel/server/internal/platform/support/logger"
	"github.com/perfect-panel/server/internal/platform/support/xerr"
	"github.com/pkg/errors"
)

type CheckVerificationCodeLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Check verification code
func NewCheckVerificationCodeLogic(ctx context.Context, deps Deps) *CheckVerificationCodeLogic {
	return &CheckVerificationCodeLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *CheckVerificationCodeLogic) CheckVerificationCode(req *types.CheckVerificationCodeRequest) (resp *types.CheckVerificationCodeRespone, err error) {
	resp = &types.CheckVerificationCodeRespone{}
	if req.Method == authmethod.Email {
		cacheKey := fmt.Sprintf("%s:%s:%s", config.AuthCodeCacheKey, config.ParseVerifyType(req.Type), req.Account)
		value, err := l.deps.Redis.Get(l.ctx, cacheKey).Result()
		if err != nil {
			return resp, nil
		}
		var payload CacheKeyPayload
		if err := json.Unmarshal([]byte(value), &payload); err != nil {
			return resp, nil
		}
		if payload.Code != req.Code {
			return resp, nil
		}
		resp.Status = true
	}
	if req.Method == authmethod.Mobile {
		if !phone.CheckPhone(req.Account) {
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.TelephoneError), "Invalid phone number")
		}
		cacheKey := fmt.Sprintf("%s:%s:+%s", config.AuthCodeTelephoneCacheKey, config.ParseVerifyType(req.Type), req.Account)
		value, err := l.deps.Redis.Get(l.ctx, cacheKey).Result()
		if err != nil {
			return resp, nil
		}
		var payload CacheKeyPayload
		if err := json.Unmarshal([]byte(value), &payload); err != nil {
			return resp, nil
		}
		if payload.Code != req.Code {
			return resp, nil
		}
		resp.Status = true
	}
	return resp, nil
}
