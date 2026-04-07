package system

import (
	"context"
	"encoding/json"
	"github.com/perfect-panel/server/initialize"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type SetNodeMultiplierInput struct {
	Body types.SetNodeMultiplierRequest
}

func SetNodeMultiplierHandler(svcCtx *svc.ServiceContext) func(context.Context, *SetNodeMultiplierInput) (*struct{}, error) {
	return func(ctx context.Context, input *SetNodeMultiplierInput) (*struct{}, error) {
		l := NewSetNodeMultiplierLogic(ctx, svcCtx)
		if err := l.SetNodeMultiplier(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type SetNodeMultiplierLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Set Node Multiplier
func NewSetNodeMultiplierLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetNodeMultiplierLogic {
	return &SetNodeMultiplierLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetNodeMultiplierLogic) SetNodeMultiplier(req *types.SetNodeMultiplierRequest) error {
	data, err := json.Marshal(req.Periods)
	if err != nil {
		l.Logger.Error("Marshal Node Multiplier Config Error: ", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "Marshal Node Multiplier Config Error: %s", err.Error())
	}
	if err = l.svcCtx.SystemModel.UpdateNodeMultiplierConfig(l.ctx, string(data)); err != nil {
		l.Logger.Error("Update Node Multiplier Config Error: ", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "Update Node Multiplier Config Error: %s", err.Error())
	}
	// update Node Multiplier
	initialize.Node(l.svcCtx)

	return nil
}
