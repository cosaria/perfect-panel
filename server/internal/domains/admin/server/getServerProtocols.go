package server

import (
	"context"

	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/internal/platform/support/logger"
	"github.com/perfect-panel/server/internal/platform/support/tool"
	"github.com/perfect-panel/server/internal/platform/support/xerr"
	"github.com/pkg/errors"
)

type GetServerProtocolsInput struct {
	types.GetServerProtocolsRequest
}

type GetServerProtocolsOutput struct {
	Body *types.GetServerProtocolsResponse
}

func GetServerProtocolsHandler(deps Deps) func(context.Context, *GetServerProtocolsInput) (*GetServerProtocolsOutput, error) {
	return func(ctx context.Context, input *GetServerProtocolsInput) (*GetServerProtocolsOutput, error) {
		l := NewGetServerProtocolsLogic(ctx, deps)
		resp, err := l.GetServerProtocols(&input.GetServerProtocolsRequest)
		if err != nil {
			return nil, err
		}
		return &GetServerProtocolsOutput{Body: resp}, nil
	}
}

type GetServerProtocolsLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Get Server Protocols
func NewGetServerProtocolsLogic(ctx context.Context, deps Deps) *GetServerProtocolsLogic {
	return &GetServerProtocolsLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *GetServerProtocolsLogic) GetServerProtocols(req *types.GetServerProtocolsRequest) (resp *types.GetServerProtocolsResponse, err error) {
	// find server
	data, err := l.deps.NodeModel.FindOneServer(l.ctx, req.Id)
	if err != nil {
		l.Errorf("[GetServerProtocols] FindOneServer Error: %s", err.Error())
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "[GetServerProtocols] FindOneServer Error: %s", err.Error())
	}

	// handler protocols
	var protocols []types.Protocol
	dst, err := data.UnmarshalProtocols()
	if err != nil {
		l.Errorf("[FilterServerList] UnmarshalProtocols Error: %s", err.Error())
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "[FilterServerList] UnmarshalProtocols Error: %s", err.Error())
	}
	tool.DeepCopy(&protocols, dst)

	return &types.GetServerProtocolsResponse{
		Protocols: protocols,
	}, nil
}
