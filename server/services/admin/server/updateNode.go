package server

import (
	"context"
	"github.com/perfect-panel/server/models/node"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/pkg/errors"
)

type UpdateNodeInput struct {
	Body types.UpdateNodeRequest
}

func UpdateNodeHandler(deps Deps) func(context.Context, *UpdateNodeInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateNodeInput) (*struct{}, error) {
		l := NewUpdateNodeLogic(ctx, deps)
		if err := l.UpdateNode(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type UpdateNodeLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// NewUpdateNodeLogic Update Node
func NewUpdateNodeLogic(ctx context.Context, deps Deps) *UpdateNodeLogic {
	return &UpdateNodeLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *UpdateNodeLogic) UpdateNode(req *types.UpdateNodeRequest) error {
	data, err := l.deps.NodeModel.FindOneNode(l.ctx, req.Id)
	if err != nil {
		l.Errorw("[UpdateNode] Query Database Error: ", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseUpdateError), "[UpdateNode] Query Database Error")
	}
	data.Name = req.Name
	data.Tags = tool.StringSliceToString(req.Tags)
	data.ServerId = req.ServerId
	data.Port = req.Port
	data.Address = req.Address
	data.Protocol = req.Protocol
	data.Enabled = req.Enabled
	err = l.deps.NodeModel.UpdateNode(l.ctx, data)
	if err != nil {
		l.Errorw("[UpdateNode] Update Database Error: ", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseUpdateError), "[UpdateNode] Update Database Error")
	}
	return l.deps.NodeModel.ClearNodeCache(l.ctx, &node.FilterNodeParams{
		Page:     1,
		Size:     1000,
		ServerId: []int64{data.ServerId},
		Search:   "",
	})
}
