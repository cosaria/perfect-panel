package server

import (
	"context"
	"github.com/perfect-panel/server/models/node"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type ResetSortWithServerInput struct {
	Body types.ResetSortRequest
}

func ResetSortWithServerHandler(svcCtx *svc.ServiceContext) func(context.Context, *ResetSortWithServerInput) (*struct{}, error) {
	return func(ctx context.Context, input *ResetSortWithServerInput) (*struct{}, error) {
		l := NewResetSortWithServerLogic(ctx, svcCtx)
		if err := l.ResetSortWithServer(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type ResetSortWithServerLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewResetSortWithServerLogic Reset server sort
func NewResetSortWithServerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetSortWithServerLogic {
	return &ResetSortWithServerLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResetSortWithServerLogic) ResetSortWithServer(req *types.ResetSortRequest) error {
	err := l.svcCtx.NodeModel.Transaction(l.ctx, func(db *gorm.DB) error {
		// find all servers id
		var existingIDs []int64
		db.Model(&node.Server{}).Select("id").Find(&existingIDs)
		// check if the id is valid
		validIDMap := make(map[int64]bool)
		for _, id := range existingIDs {
			validIDMap[id] = true
		}
		// check if the sort is valid
		var validItems []types.SortItem
		for _, item := range req.Sort {
			if validIDMap[item.Id] {
				validItems = append(validItems, item)
			}
		}
		// query all servers
		var servers []*node.Server
		db.Model(&node.Server{}).Order("sort ASC").Find(&servers)
		// create a map of the current sort
		currentSortMap := make(map[int64]int64)
		for _, item := range servers {
			currentSortMap[item.Id] = int64(item.Sort)
		}

		// new sort map
		newSortMap := make(map[int64]int64)
		for _, item := range validItems {
			newSortMap[item.Id] = item.Sort
		}

		var itemsToUpdate []types.SortItem
		for _, item := range validItems {
			if oldSort, exists := currentSortMap[item.Id]; exists && oldSort != item.Sort {
				itemsToUpdate = append(itemsToUpdate, item)
			}
		}
		for _, item := range itemsToUpdate {
			s, err := l.svcCtx.NodeModel.FindOneServer(l.ctx, item.Id)
			if err != nil {
				return err
			}
			s.Sort = int(item.Sort)
			if err = l.svcCtx.NodeModel.UpdateServer(l.ctx, s, db); err != nil {
				l.Errorw("[NodeSort] Update Database Error: ", logger.Field("error", err.Error()), logger.Field("id", item.Id), logger.Field("sort", item.Sort))
				return err
			}
		}
		return nil
	})
	if err != nil {
		l.Errorw("[NodeSort] Update Database Error: ", logger.Field("error", err.Error()))
		return errors.Wrap(xerr.NewErrCode(xerr.DatabaseUpdateError), err.Error())
	}
	return nil
}
