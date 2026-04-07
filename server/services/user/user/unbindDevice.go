package user

import (
	"context"
	"fmt"

	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/models/user"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type UnbindDeviceInput struct {
	Body types.UnbindDeviceRequest
}

func UnbindDeviceHandler(deps Deps) func(context.Context, *UnbindDeviceInput) (*struct{}, error) {
	return func(ctx context.Context, input *UnbindDeviceInput) (*struct{}, error) {
		l := NewUnbindDeviceLogic(ctx, deps)
		if err := l.UnbindDevice(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type UnbindDeviceLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Unbind Device
func NewUnbindDeviceLogic(ctx context.Context, deps Deps) *UnbindDeviceLogic {
	return &UnbindDeviceLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *UnbindDeviceLogic) UnbindDevice(req *types.UnbindDeviceRequest) error {
	userInfo := l.ctx.Value(config.CtxKeyUser).(*user.User)
	device, err := l.deps.UserModel.FindOneDevice(l.ctx, req.Id)
	if err != nil {
		return errors.Wrapf(xerr.NewErrCode(xerr.DeviceNotExist), "find device")
	}

	if device.UserId != userInfo.Id {
		return errors.Wrapf(xerr.NewErrCode(xerr.InvalidParams), "device not belong to user")
	}

	return l.deps.DB.Transaction(func(tx *gorm.DB) error {
		var deleteDevice user.Device
		err = tx.Model(&deleteDevice).Where("id = ?", req.Id).First(&deleteDevice).Error
		if err != nil {
			return errors.Wrapf(xerr.NewErrCode(xerr.QueueEnqueueError), "find device err: %v", err)
		}
		err = tx.Delete(deleteDevice).Error
		if err != nil {
			return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseDeletedError), "delete device err: %v", err)
		}
		var userAuth user.AuthMethods
		err = tx.Model(&userAuth).Where("auth_identifier = ? and auth_type = ?", deleteDevice.Identifier, "device").First(&userAuth).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}
			return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "find device online record err: %v", err)
		}

		err = tx.Delete(&userAuth).Error
		if err != nil {
			return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseDeletedError), "delete device online record err: %v", err)
		}
		sessionId := l.ctx.Value(config.CtxKeySessionID)
		sessionIdCacheKey := fmt.Sprintf("%v:%v", config.SessionIdKey, sessionId)
		l.deps.Redis.Del(l.ctx, sessionIdCacheKey)
		return nil
	})
}
