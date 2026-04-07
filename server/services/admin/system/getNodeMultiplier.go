package system

import (
	"context"
	"encoding/json"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type GetNodeMultiplierOutput struct {
	Body *types.GetNodeMultiplierResponse
}

func GetNodeMultiplierHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*GetNodeMultiplierOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetNodeMultiplierOutput, error) {
		l := NewGetNodeMultiplierLogic(ctx, svcCtx)
		resp, err := l.GetNodeMultiplier()
		if err != nil {
			return nil, err
		}
		return &GetNodeMultiplierOutput{Body: resp}, nil
	}
}

type GetNodeMultiplierLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Get Node Multiplier
func NewGetNodeMultiplierLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetNodeMultiplierLogic {
	return &GetNodeMultiplierLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetNodeMultiplierLogic) GetNodeMultiplier() (resp *types.GetNodeMultiplierResponse, err error) {
	data, err := l.svcCtx.SystemModel.FindNodeMultiplierConfig(l.ctx)
	if err != nil {
		l.Logger.Error("Get Node Multiplier Config Error: ", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "Get Node Multiplier Config Error: %s", err.Error())
	}
	var periods []types.TimePeriod
	if data.Value != "" {
		if err := json.Unmarshal([]byte(data.Value), &periods); err != nil {
			l.Logger.Error("Unmarshal Node Multiplier Config Error: ", logger.Field("error", err.Error()), logger.Field("value", data.Value))
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "Unmarshal Node Multiplier Config Error: %s", err.Error())
		}
	}

	return &types.GetNodeMultiplierResponse{
		Periods: periods,
	}, nil
}
