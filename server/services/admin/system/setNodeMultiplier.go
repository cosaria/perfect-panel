package system

import (
	"context"
	"encoding/json"

	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/pkg/errors"
)

type SetNodeMultiplierInput struct {
	Body types.SetNodeMultiplierRequest
}

func SetNodeMultiplierHandler(deps Deps) func(context.Context, *SetNodeMultiplierInput) (*struct{}, error) {
	return func(ctx context.Context, input *SetNodeMultiplierInput) (*struct{}, error) {
		l := NewSetNodeMultiplierLogic(ctx, deps)
		if err := l.SetNodeMultiplier(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type SetNodeMultiplierLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Set Node Multiplier
func NewSetNodeMultiplierLogic(ctx context.Context, deps Deps) *SetNodeMultiplierLogic {
	return &SetNodeMultiplierLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *SetNodeMultiplierLogic) SetNodeMultiplier(req *types.SetNodeMultiplierRequest) error {
	data, err := json.Marshal(req.Periods)
	if err != nil {
		l.Error("Marshal Node Multiplier Config Error: ", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "Marshal Node Multiplier Config Error: %s", err.Error())
	}
	if err = l.deps.SystemModel.UpdateNodeMultiplierConfig(l.ctx, string(data)); err != nil {
		l.Error("Update Node Multiplier Config Error: ", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "Update Node Multiplier Config Error: %s", err.Error())
	}
	if l.deps.ReloadNode != nil {
		l.deps.ReloadNode()
	}

	return nil
}
