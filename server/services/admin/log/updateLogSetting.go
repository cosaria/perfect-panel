package log

import (
	"context"
	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/models/system"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"reflect"
)

type UpdateLogSettingInput struct {
	Body types.LogSetting
}

func UpdateLogSettingHandler(svcCtx *svc.ServiceContext) func(context.Context, *UpdateLogSettingInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateLogSettingInput) (*struct{}, error) {
		l := NewUpdateLogSettingLogic(ctx, svcCtx)
		if err := l.UpdateLogSetting(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type UpdateLogSettingLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewUpdateLogSettingLogic Update log setting
func NewUpdateLogSettingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogSettingLogic {
	return &UpdateLogSettingLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLogSettingLogic) UpdateLogSetting(req *types.LogSetting) error {
	v := reflect.ValueOf(*req)
	// Get the reflection type of the structure
	t := v.Type()
	err := l.svcCtx.SystemModel.Transaction(l.ctx, func(db *gorm.DB) error {
		var err error
		for i := 0; i < v.NumField(); i++ {
			// Get the field name
			fieldName := t.Field(i).Name
			// Get the field value to string
			fieldValue := tool.ConvertValueToString(v.Field(i))
			// Update the server config
			err = db.Model(&system.System{}).Where("`category` = 'log' and `key` = ?", fieldName).Update("value", fieldValue).Error
			if err != nil {
				break
			}
		}
		return err
	})
	if err != nil {
		l.Errorw("[UpdateLogSetting] update log setting error", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseUpdateError), " update log setting error: %v", err)
	}

	l.svcCtx.Config.Log = config.Log{
		AutoClear: *req.AutoClear,
		ClearDays: req.ClearDays,
	}

	return nil
}
