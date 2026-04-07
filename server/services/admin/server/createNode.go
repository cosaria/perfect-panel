package server

import (
	"context"
	"github.com/perfect-panel/server/models/node"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type CreateNodeInput struct {
	Body types.CreateNodeRequest
}

func CreateNodeHandler(svcCtx *svc.ServiceContext) func(context.Context, *CreateNodeInput) (*struct{}, error) {
	return func(ctx context.Context, input *CreateNodeInput) (*struct{}, error) {
		l := NewCreateNodeLogic(ctx, svcCtx)
		if err := l.CreateNode(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type CreateNodeLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewCreateNodeLogic Create Node
func NewCreateNodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateNodeLogic {
	return &CreateNodeLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
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
	err := l.svcCtx.NodeModel.InsertNode(l.ctx, &data)
	if err != nil {
		l.Errorw("[CreateNode] Insert Database Error: ", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseInsertError), "[CreateNode] Insert Database Error")
	}

	return nil
}
