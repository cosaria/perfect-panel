package server

import (
	"context"
	"github.com/perfect-panel/server/models/node"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
	"strings"
)

type ToggleNodeStatusInput struct {
	Body types.ToggleNodeStatusRequest
}

func ToggleNodeStatusHandler(svcCtx *svc.ServiceContext) func(context.Context, *ToggleNodeStatusInput) (*struct{}, error) {
	return func(ctx context.Context, input *ToggleNodeStatusInput) (*struct{}, error) {
		l := NewToggleNodeStatusLogic(ctx, svcCtx)
		if err := l.ToggleNodeStatus(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type ToggleNodeStatusLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewToggleNodeStatusLogic Toggle Node Status
func NewToggleNodeStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ToggleNodeStatusLogic {
	return &ToggleNodeStatusLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ToggleNodeStatusLogic) ToggleNodeStatus(req *types.ToggleNodeStatusRequest) error {
	data, err := l.svcCtx.NodeModel.FindOneNode(l.ctx, req.Id)
	if err != nil {
		l.Errorw("[ToggleNodeStatus] Query Database Error: ", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "[ToggleNodeStatus] Query Database Error")
	}
	data.Enabled = req.Enable

	err = l.svcCtx.NodeModel.UpdateNode(l.ctx, data)
	if err != nil {
		l.Errorw("[ToggleNodeStatus] Update Database Error: ", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseUpdateError), "[ToggleNodeStatus] Update Database Error")
	}

	return l.svcCtx.NodeModel.ClearNodeCache(l.ctx, &node.FilterNodeParams{
		Page:     1,
		Size:     1000,
		ServerId: []int64{data.ServerId},
		Tag:      strings.Split(data.Tags, ","),
		Search:   "",
	})
}
