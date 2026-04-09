package system

import (
	"context"
	"encoding/json"

	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/internal/platform/support/logger"
	"github.com/perfect-panel/server/internal/platform/support/xerr"
	"github.com/pkg/errors"
)

type GetNodeMultiplierOutput struct {
	Body *types.GetNodeMultiplierResponse
}

func GetNodeMultiplierHandler(deps Deps) func(context.Context, *struct{}) (*GetNodeMultiplierOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetNodeMultiplierOutput, error) {
		l := NewGetNodeMultiplierLogic(ctx, deps)
		resp, err := l.GetNodeMultiplier()
		if err != nil {
			return nil, err
		}
		return &GetNodeMultiplierOutput{Body: resp}, nil
	}
}

type GetNodeMultiplierLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Get Node Multiplier
func NewGetNodeMultiplierLogic(ctx context.Context, deps Deps) *GetNodeMultiplierLogic {
	return &GetNodeMultiplierLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *GetNodeMultiplierLogic) GetNodeMultiplier() (resp *types.GetNodeMultiplierResponse, err error) {
	data, err := l.deps.SystemModel.FindNodeMultiplierConfig(l.ctx)
	if err != nil {
		l.Error("Get Node Multiplier Config Error: ", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "Get Node Multiplier Config Error: %s", err.Error())
	}
	var periods []types.TimePeriod
	if data.Value != "" {
		if err := json.Unmarshal([]byte(data.Value), &periods); err != nil {
			l.Error("Unmarshal Node Multiplier Config Error: ", logger.Field("error", err.Error()), logger.Field("value", data.Value))
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "Unmarshal Node Multiplier Config Error: %s", err.Error())
		}
	}

	return &types.GetNodeMultiplierResponse{
		Periods: periods,
	}, nil
}
