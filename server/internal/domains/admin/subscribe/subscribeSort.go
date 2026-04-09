package subscribe

import (
	"context"

	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/models/subscribe"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type SubscribeSortInput struct {
	Body types.SubscribeSortRequest
}

func SubscribeSortHandler(deps Deps) func(context.Context, *SubscribeSortInput) (*struct{}, error) {
	return func(ctx context.Context, input *SubscribeSortInput) (*struct{}, error) {
		l := NewSubscribeSortLogic(ctx, deps)
		if err := l.SubscribeSort(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type SubscribeSortLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// NewSubscribeSortLogic Subscribe sort
func NewSubscribeSortLogic(ctx context.Context, deps Deps) *SubscribeSortLogic {
	return &SubscribeSortLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *SubscribeSortLogic) SubscribeSort(req *types.SubscribeSortRequest) error {
	var sort = make(map[int64]int64, len(req.Sort))
	var ids []int64
	for i, v := range req.Sort {
		sort[v.Id] = int64(i)
		ids = append(ids, v.Id)
	}
	// query min sort by ids
	minSort, err := l.deps.SubscribeModel.QuerySubscribeMinSortByIds(l.ctx, ids)
	if err != nil {
		l.Error("[SubscribeSortLogic] query subscribe list by ids error: ", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "query subscribe list by ids error: %v", err.Error())
	}
	_, subs, err := l.deps.SubscribeModel.FilterList(l.ctx, &subscribe.FilterParams{
		Page: 1,
		Size: 9999,
		Ids:  ids,
	})
	if err != nil {
		l.Error("[SubscribeSortLogic] query subscribe list by ids error: ", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "query subscribe list by ids error: %v", err.Error())
	}
	// reordering
	for _, sub := range subs {
		if newSort, ok := sort[sub.Id]; ok {
			sub.Sort = minSort + newSort
		}
	}
	// update sort
	err = l.deps.SubscribeModel.Transaction(l.ctx, func(db *gorm.DB) error {
		return db.Save(subs).Error
	})
	if err != nil {
		l.Error("[SubscribeSortLogic] update subscribe sort error: ", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseUpdateError), "update subscribe sort error: %v", err.Error())
	}
	l.Info("[UpdateSubscribeSort] Successfully updated subscribe sort")
	return nil
}
