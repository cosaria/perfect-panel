package server

import (
	"context"
	"strings"

	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/models/node"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/pkg/errors"
)

type ToggleNodeStatusInput struct {
	Body types.ToggleNodeStatusRequest
}

func ToggleNodeStatusHandler(deps Deps) func(context.Context, *ToggleNodeStatusInput) (*struct{}, error) {
	return func(ctx context.Context, input *ToggleNodeStatusInput) (*struct{}, error) {
		l := NewToggleNodeStatusLogic(ctx, deps)
		if err := l.ToggleNodeStatus(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type ToggleNodeStatusLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// NewToggleNodeStatusLogic Toggle Node Status
func NewToggleNodeStatusLogic(ctx context.Context, deps Deps) *ToggleNodeStatusLogic {
	return &ToggleNodeStatusLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *ToggleNodeStatusLogic) ToggleNodeStatus(req *types.ToggleNodeStatusRequest) error {
	data, err := l.deps.NodeModel.FindOneNode(l.ctx, req.Id)
	if err != nil {
		l.Errorw("[ToggleNodeStatus] Query Database Error: ", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "[ToggleNodeStatus] Query Database Error")
	}
	data.Enabled = req.Enable

	err = l.deps.NodeModel.UpdateNode(l.ctx, data)
	if err != nil {
		l.Errorw("[ToggleNodeStatus] Update Database Error: ", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseUpdateError), "[ToggleNodeStatus] Update Database Error")
	}

	return l.deps.NodeModel.ClearNodeCache(l.ctx, &node.FilterNodeParams{
		Page:     1,
		Size:     1000,
		ServerId: []int64{data.ServerId},
		Tag:      strings.Split(data.Tags, ","),
		Search:   "",
	})
}
