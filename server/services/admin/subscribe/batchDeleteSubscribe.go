package subscribe

import (
	"context"
	"github.com/perfect-panel/server/models/user"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type BatchDeleteSubscribeInput struct {
	Body types.BatchDeleteSubscribeRequest
}

func BatchDeleteSubscribeHandler(deps Deps) func(context.Context, *BatchDeleteSubscribeInput) (*struct{}, error) {
	return func(ctx context.Context, input *BatchDeleteSubscribeInput) (*struct{}, error) {
		l := NewBatchDeleteSubscribeLogic(ctx, deps)
		if err := l.BatchDeleteSubscribe(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type BatchDeleteSubscribeLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Batch delete subscribe
func NewBatchDeleteSubscribeLogic(ctx context.Context, deps Deps) *BatchDeleteSubscribeLogic {
	return &BatchDeleteSubscribeLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

var errorIsExistActiveUser = errors.New("subscription ID belongs to an active user subscription")

func (l *BatchDeleteSubscribeLogic) BatchDeleteSubscribe(req *types.BatchDeleteSubscribeRequest) error {
	err := l.deps.DB.Transaction(func(tx *gorm.DB) error {
		for _, id := range req.Ids {
			var count int64
			// Validate whether the subscription ID belongs to an active user subscription.
			if err := tx.Model(&user.Subscribe{}).Where("subscribe_id = ? AND status = 1", id).Count(&count).Find(&user.Subscribe{}).Error; err != nil {
				l.Error("[BatchDeleteSubscribe] Query Subscribe Error: ", logger.Field("error", err.Error()))
				return err
			}
			if count > 0 {
				return errorIsExistActiveUser
			}
			if err := l.deps.SubscribeModel.Delete(l.ctx, id, tx); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, errorIsExistActiveUser) {
			return errors.Wrapf(xerr.NewErrCode(xerr.SubscribeIsUsedError), "subscription ID belongs to an active user subscription")
		}
		l.Error("[BatchDeleteSubscribe] Transaction Error: ", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseDeletedError), "delete subscribe failed: %v", err.Error())
	}
	return nil
}
