package server

import (
	"context"

	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/models/node"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/pkg/errors"
)

type DeleteServerInput struct {
	Body types.DeleteServerRequest
}

func DeleteServerHandler(deps Deps) func(context.Context, *DeleteServerInput) (*struct{}, error) {
	return func(ctx context.Context, input *DeleteServerInput) (*struct{}, error) {
		l := NewDeleteServerLogic(ctx, deps)
		if err := l.DeleteServer(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type DeleteServerLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// NewDeleteServerLogic Delete Server
func NewDeleteServerLogic(ctx context.Context, deps Deps) *DeleteServerLogic {
	return &DeleteServerLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *DeleteServerLogic) DeleteServer(req *types.DeleteServerRequest) error {
	err := l.deps.NodeModel.DeleteServer(l.ctx, req.Id)
	if err != nil {
		l.Errorw("[DeleteServer] Delete Server Error: ", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseDeletedError), "[DeleteServer] Delete Server Error")
	}
	return l.deps.NodeModel.ClearNodeCache(l.ctx, &node.FilterNodeParams{
		Page:     1,
		Size:     1000,
		ServerId: []int64{req.Id},
		Search:   "",
	})
}
