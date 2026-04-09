package server

import (
	"context"
	"strings"

	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/internal/platform/persistence/node"
	"github.com/perfect-panel/server/internal/platform/support/logger"
	"github.com/perfect-panel/server/internal/platform/support/xerr"
	"github.com/pkg/errors"
)

type DeleteNodeInput struct {
	Body types.DeleteNodeRequest
}

func DeleteNodeHandler(deps Deps) func(context.Context, *DeleteNodeInput) (*struct{}, error) {
	return func(ctx context.Context, input *DeleteNodeInput) (*struct{}, error) {
		l := NewDeleteNodeLogic(ctx, deps)
		if err := l.DeleteNode(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type DeleteNodeLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// NewDeleteNodeLogic Delete Node
func NewDeleteNodeLogic(ctx context.Context, deps Deps) *DeleteNodeLogic {
	return &DeleteNodeLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *DeleteNodeLogic) DeleteNode(req *types.DeleteNodeRequest) error {
	data, err := l.deps.NodeModel.FindOneNode(l.ctx, req.Id)
	if err != nil {
		return err
	}

	err = l.deps.NodeModel.DeleteNode(l.ctx, req.Id)
	if err != nil {
		l.Errorw("[DeleteNode] Delete Database Error: ", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseDeletedError), "[DeleteNode] Delete Database Error")
	}

	return l.deps.NodeModel.ClearNodeCache(l.ctx, &node.FilterNodeParams{
		Page:     1,
		Size:     1000,
		ServerId: []int64{data.ServerId},
		Tag:      strings.Split(data.Tags, ","),
		Search:   "",
		Protocol: data.Protocol,
	})
}
