package subscribe

import (
	"context"
	"strconv"
	"time"

	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/internal/platform/persistence/user"
	"github.com/perfect-panel/server/internal/platform/support/logger"
	"github.com/perfect-panel/server/internal/platform/support/uuidx"
	"github.com/perfect-panel/server/internal/platform/support/xerr"
	"github.com/pkg/errors"
)

type ResetAllSubscribeTokenOutput struct {
	Body *types.ResetAllSubscribeTokenResponse
}

func ResetAllSubscribeTokenHandler(deps Deps) func(context.Context, *struct{}) (*ResetAllSubscribeTokenOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*ResetAllSubscribeTokenOutput, error) {
		l := NewResetAllSubscribeTokenLogic(ctx, deps)
		resp, err := l.ResetAllSubscribeToken()
		if err != nil {
			return nil, err
		}
		return &ResetAllSubscribeTokenOutput{Body: resp}, nil
	}
}

type ResetAllSubscribeTokenLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Reset all subscribe tokens
func NewResetAllSubscribeTokenLogic(ctx context.Context, deps Deps) *ResetAllSubscribeTokenLogic {
	return &ResetAllSubscribeTokenLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *ResetAllSubscribeTokenLogic) ResetAllSubscribeToken() (resp *types.ResetAllSubscribeTokenResponse, err error) {
	var list []*user.Subscribe
	tx := l.deps.DB.WithContext(l.ctx).Begin()
	// select all active and Finished subscriptions
	if err = tx.Model(&user.Subscribe{}).Where("`status` IN ?", []int64{1, 2}).Find(&list).Error; err != nil {
		logger.Errorf("[ResetAllSubscribeToken] Failed to fetch subscribe list: %v", err.Error())
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "Failed to fetch subscribe list: %v", err.Error())
	}

	for _, sub := range list {
		sub.Token = uuidx.SubscribeToken(strconv.FormatInt(time.Now().UnixMilli(), 10) + strconv.FormatInt(sub.Id, 10))
		sub.UUID = uuidx.NewUUID().String()
		if err = tx.Model(&user.Subscribe{}).Where("id = ?", sub.Id).Save(sub).Error; err != nil {
			tx.Rollback()
			logger.Errorf("[ResetAllSubscribeToken] Failed to update subscribe token for ID %d: %v", sub.Id, err.Error())
			return &types.ResetAllSubscribeTokenResponse{
				Success: false,
			}, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseUpdateError), "Failed to update subscribe token for ID %d: %v", sub.Id, err.Error())
		}
	}
	if err = tx.Commit().Error; err != nil {
		logger.Errorf("[ResetAllSubscribeToken] Failed to commit transaction: %v", err.Error())
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseUpdateError), "Failed to commit transaction: %v", err.Error())
	}

	return &types.ResetAllSubscribeTokenResponse{
		Success: true,
	}, nil
}
