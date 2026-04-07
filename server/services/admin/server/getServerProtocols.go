package server

import (
	"context"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type GetServerProtocolsInput struct {
	types.GetServerProtocolsRequest
}

type GetServerProtocolsOutput struct {
	Body *types.GetServerProtocolsResponse
}

func GetServerProtocolsHandler(svcCtx *svc.ServiceContext) func(context.Context, *GetServerProtocolsInput) (*GetServerProtocolsOutput, error) {
	return func(ctx context.Context, input *GetServerProtocolsInput) (*GetServerProtocolsOutput, error) {
		l := NewGetServerProtocolsLogic(ctx, svcCtx)
		resp, err := l.GetServerProtocols(&input.GetServerProtocolsRequest)
		if err != nil {
			return nil, err
		}
		return &GetServerProtocolsOutput{Body: resp}, nil
	}
}

type GetServerProtocolsLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Get Server Protocols
func NewGetServerProtocolsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetServerProtocolsLogic {
	return &GetServerProtocolsLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetServerProtocolsLogic) GetServerProtocols(req *types.GetServerProtocolsRequest) (resp *types.GetServerProtocolsResponse, err error) {
	// find server
	data, err := l.svcCtx.NodeModel.FindOneServer(l.ctx, req.Id)
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
