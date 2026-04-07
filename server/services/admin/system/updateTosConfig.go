package system

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

type UpdateTosConfigInput struct {
	Body types.TosConfig
}

func UpdateTosConfigHandler(svcCtx *svc.ServiceContext) func(context.Context, *UpdateTosConfigInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateTosConfigInput) (*struct{}, error) {
		l := NewUpdateTosConfigLogic(ctx, svcCtx)
		if err := l.UpdateTosConfig(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type UpdateTosConfigLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateTosConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTosConfigLogic {
	return &UpdateTosConfigLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateTosConfigLogic) UpdateTosConfig(req *types.TosConfig) error {
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
			// Update the tos config
			err = db.Model(&system.System{}).Where("`category` = 'tos' and `key` = ?", fieldName).Update("value", fieldValue).Error
			if err != nil {
				break
			}
		}
		return l.svcCtx.Redis.Del(l.ctx, config.TosConfigKey).Err()
	})
	if err != nil {
		l.Errorw("[UpdateTosConfigLogic] update tos config error: ", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseUpdateError), "update tos config error: %v", err)
	}

	return nil
}
