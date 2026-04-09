package subscribe

import (
	"context"

	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/models/user"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type DeleteSubscribeInput struct {
	Body types.DeleteSubscribeRequest
}

func DeleteSubscribeHandler(deps Deps) func(context.Context, *DeleteSubscribeInput) (*struct{}, error) {
	return func(ctx context.Context, input *DeleteSubscribeInput) (*struct{}, error) {
		l := NewDeleteSubscribeLogic(ctx, deps)
		if err := l.DeleteSubscribe(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type DeleteSubscribeLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Delete subscribe
func NewDeleteSubscribeLogic(ctx context.Context, deps Deps) *DeleteSubscribeLogic {
	return &DeleteSubscribeLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *DeleteSubscribeLogic) DeleteSubscribe(req *types.DeleteSubscribeRequest) error {
	// Check if the subscribe exists
	var total int64
	err := l.deps.UserModel.Transaction(l.ctx, func(db *gorm.DB) error {
		return db.Model(&user.Subscribe{}).Where("subscribe_id = ? AND `status` = ?", req.Id, 1).Count(&total).Find(&user.Subscribe{}).Error
	})
	if err != nil {
		l.Error("[DeleteSubscribeLogic] check subscribe failed: ", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "check subscribe failed: %v", err.Error())
	}
	if total != 0 {
		return errors.Wrapf(xerr.NewErrCode(xerr.SubscribeIsUsedError), "subscribe is used")
	}

	err = l.deps.SubscribeModel.Delete(l.ctx, req.Id)
	if err != nil {
		l.Error("[DeleteSubscribeLogic] delete subscribe failed: ", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseDeletedError), "delete subscribe failed: %v", err.Error())
	}
	return nil
}
