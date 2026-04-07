package system

import (
	"context"
	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/initialize"
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

type UpdateVerifyConfigInput struct {
	Body types.VerifyConfig
}

func UpdateVerifyConfigHandler(svcCtx *svc.ServiceContext) func(context.Context, *UpdateVerifyConfigInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateVerifyConfigInput) (*struct{}, error) {
		l := NewUpdateVerifyConfigLogic(ctx, svcCtx)
		if err := l.UpdateVerifyConfig(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type UpdateVerifyConfigLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateVerifyConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateVerifyConfigLogic {
	return &UpdateVerifyConfigLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateVerifyConfigLogic) UpdateVerifyConfig(req *types.VerifyConfig) error {
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
			// Update the site config
			err = db.Model(&system.System{}).Where("`category` = 'verify' and `key` = ?", fieldName).Update("value", fieldValue).Error
			if err != nil {
				break
			}
		}
		if err != nil {
			return err
		}
		// clear cache
		return l.svcCtx.Redis.Del(l.ctx, config.VerifyConfigKey, config.GlobalConfigKey).Err()
	})
	if err != nil {
		l.Errorw("[UpdateVerifyConfigLogic] update verify config error: ", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseUpdateError), "update verify config error: %v", err)
	}
	// Update the config
	tool.DeepCopy(&l.svcCtx.Config.Verify, req)
	initialize.Verify(l.svcCtx)
	return nil
}
