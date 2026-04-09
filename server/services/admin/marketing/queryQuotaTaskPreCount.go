package marketing

import (
	"context"
	"github.com/perfect-panel/server/models/user"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/internal/platform/http/types"
	"time"
)

type QueryQuotaTaskPreCountInput struct {
	Body types.QueryQuotaTaskPreCountRequest
}

type QueryQuotaTaskPreCountOutput struct {
	Body *types.QueryQuotaTaskPreCountResponse
}

func QueryQuotaTaskPreCountHandler(deps Deps) func(context.Context, *QueryQuotaTaskPreCountInput) (*QueryQuotaTaskPreCountOutput, error) {
	return func(ctx context.Context, input *QueryQuotaTaskPreCountInput) (*QueryQuotaTaskPreCountOutput, error) {
		l := NewQueryQuotaTaskPreCountLogic(ctx, deps)
		resp, err := l.QueryQuotaTaskPreCount(&input.Body)
		if err != nil {
			return nil, err
		}
		return &QueryQuotaTaskPreCountOutput{Body: resp}, nil
	}
}

type QueryQuotaTaskPreCountLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// NewQueryQuotaTaskPreCountLogic Query quota task pre-count
func NewQueryQuotaTaskPreCountLogic(ctx context.Context, deps Deps) *QueryQuotaTaskPreCountLogic {
	return &QueryQuotaTaskPreCountLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *QueryQuotaTaskPreCountLogic) QueryQuotaTaskPreCount(req *types.QueryQuotaTaskPreCountRequest) (resp *types.QueryQuotaTaskPreCountResponse, err error) {
	tx := l.deps.DB.WithContext(l.ctx).Model(&user.Subscribe{})
	var count int64

	if len(req.Subscribers) > 0 {
		tx = tx.Where("`subscribe_id` IN ?", req.Subscribers)
	}

	if req.IsActive != nil && *req.IsActive {
		tx = tx.Where("`status` IN ?", []int64{0, 1, 2}) // 0: Pending 1: Active 2: Finished
	}
	if req.StartTime != 0 {
		start := time.UnixMilli(req.StartTime)
		tx = tx.Where("`start_time` <= ?", start)
	}
	if req.EndTime != 0 {
		end := time.UnixMilli(req.EndTime)
		tx = tx.Where("`expire_time` >= ?", end)
	}
	if err = tx.Count(&count).Error; err != nil {
		l.Errorf("[QueryQuotaTaskPreCount] count error: %v", err.Error())
		return nil, err
	}

	return &types.QueryQuotaTaskPreCountResponse{
		Count: count,
	}, nil
}
