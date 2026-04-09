package marketing

import (
	"context"
	"github.com/perfect-panel/server/models/task"
	"github.com/perfect-panel/server/models/user"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/internal/platform/http/types"
	"gorm.io/gorm"
	"time"
)

type GetPreSendEmailCountInput struct {
	Body types.GetPreSendEmailCountRequest
}

type GetPreSendEmailCountOutput struct {
	Body *types.GetPreSendEmailCountResponse
}

func GetPreSendEmailCountHandler(deps Deps) func(context.Context, *GetPreSendEmailCountInput) (*GetPreSendEmailCountOutput, error) {
	return func(ctx context.Context, input *GetPreSendEmailCountInput) (*GetPreSendEmailCountOutput, error) {
		l := NewGetPreSendEmailCountLogic(ctx, deps)
		resp, err := l.GetPreSendEmailCount(&input.Body)
		if err != nil {
			return nil, err
		}
		return &GetPreSendEmailCountOutput{Body: resp}, nil
	}
}

type GetPreSendEmailCountLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// NewGetPreSendEmailCountLogic Get pre-send email count
func NewGetPreSendEmailCountLogic(ctx context.Context, deps Deps) *GetPreSendEmailCountLogic {
	return &GetPreSendEmailCountLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *GetPreSendEmailCountLogic) GetPreSendEmailCount(req *types.GetPreSendEmailCountRequest) (resp *types.GetPreSendEmailCountResponse, err error) {
	tx := l.deps.DB
	var count int64
	// 通用查询器（含 user JOIN + 注册时间范围过滤）
	baseQuery := func() *gorm.DB {
		query := tx.Model(&user.AuthMethods{}).
			Select("auth_identifier").
			Joins("JOIN user ON user.id = user_auth_methods.user_id").
			Where("auth_type = ?", "email")

		if req.RegisterStartTime != 0 {

			registerStartTime := time.UnixMilli(req.RegisterStartTime)

			query = query.Where("user.created_at >= ?", registerStartTime)
		}
		if req.RegisterEndTime != 0 {
			registerEndTime := time.UnixMilli(req.RegisterEndTime)
			query = query.Where("user.created_at <= ?", registerEndTime)
		}
		return query
	}
	var query *gorm.DB
	scope := task.ParseScopeType(req.Scope)

	switch scope {
	case task.ScopeAll:
		query = baseQuery()

	case task.ScopeActive:
		query = baseQuery().
			Joins("JOIN user_subscribe ON user.id = user_subscribe.user_id").
			Where("user_subscribe.status IN ?", []int64{1, 2})

	case task.ScopeExpired:
		query = baseQuery().
			Joins("JOIN user_subscribe ON user.id = user_subscribe.user_id").
			Where("user_subscribe.status = ?", 3)

	case task.ScopeNone:
		query = baseQuery().
			Joins("LEFT JOIN user_subscribe ON user.id = user_subscribe.user_id").
			Where("user_subscribe.user_id IS NULL")
	case task.ScopeSkip:
		// Skip scope does not require a count
		query = nil
	default:
		l.Errorf("[CreateBatchSendEmailTask] Invalid scope: %v", req.Scope)
		return nil, xerr.NewErrMsg("Invalid email scope")

	}

	if query != nil {
		if err = query.Count(&count).Error; err != nil {
			l.Errorf("[GetPreSendEmailCount] Count error: %v", err)
			return nil, xerr.NewErrMsg("Failed to count emails")
		}
	}

	return &types.GetPreSendEmailCountResponse{
		Count: count,
	}, nil
}
