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

type CreateNodeInput struct {
	Body types.CreateNodeRequest
}

func CreateNodeHandler(deps Deps) func(context.Context, *CreateNodeInput) (*struct{}, error) {
	return func(ctx context.Context, input *CreateNodeInput) (*struct{}, error) {
		l := NewCreateNodeLogic(ctx, deps)
		if err := l.CreateNode(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type CreateNodeLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// NewCreateNodeLogic Create Node
func NewCreateNodeLogic(ctx context.Context, deps Deps) *CreateNodeLogic {
	return &CreateNodeLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *CreateNodeLogic) CreateNode(req *types.CreateNodeRequest) error {
	data := node.Node{
		Name:     req.Name,
		Tags:     tool.StringSliceToString(req.Tags),
		Enabled:  req.Enabled,
		Port:     req.Port,
		Address:  req.Address,
		ServerId: req.ServerId,
		Protocol: req.Protocol,
	}
	err := l.deps.NodeModel.InsertNode(l.ctx, &data)
	if err != nil {
		l.Errorw("[CreateNode] Insert Database Error: ", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseInsertError), "[CreateNode] Insert Database Error")
	}

	return nil
}
