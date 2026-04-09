package subscribe

import (
	"context"

	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/internal/platform/persistence/subscribe"
	"github.com/perfect-panel/server/internal/platform/support/logger"
	"github.com/perfect-panel/server/internal/platform/support/xerr"
	"github.com/pkg/errors"
)

type BatchDeleteSubscribeGroupInput struct {
	Body types.BatchDeleteSubscribeGroupRequest
}

func BatchDeleteSubscribeGroupHandler(deps Deps) func(context.Context, *BatchDeleteSubscribeGroupInput) (*struct{}, error) {
	return func(ctx context.Context, input *BatchDeleteSubscribeGroupInput) (*struct{}, error) {
		l := NewBatchDeleteSubscribeGroupLogic(ctx, deps)
		if err := l.BatchDeleteSubscribeGroup(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type BatchDeleteSubscribeGroupLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Batch delete subscribe group
func NewBatchDeleteSubscribeGroupLogic(ctx context.Context, deps Deps) *BatchDeleteSubscribeGroupLogic {
	return &BatchDeleteSubscribeGroupLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *BatchDeleteSubscribeGroupLogic) BatchDeleteSubscribeGroup(req *types.BatchDeleteSubscribeGroupRequest) error {
	err := l.deps.DB.Model(&subscribe.Group{}).Where("id in ?", req.Ids).Delete(&subscribe.Group{}).Error
	if err != nil {
		l.Error("[BatchDeleteSubscribeGroup] Delete Database Error: ", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseDeletedError), "batch delete subscribe group failed: %v", err.Error())
	}
	return nil
}
