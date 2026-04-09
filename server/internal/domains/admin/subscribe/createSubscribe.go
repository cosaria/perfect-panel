package subscribe

import (
	"context"
	"encoding/json"

	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/models/subscribe"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/pkg/errors"
)

type CreateSubscribeInput struct {
	Body types.CreateSubscribeRequest
}

func CreateSubscribeHandler(deps Deps) func(context.Context, *CreateSubscribeInput) (*struct{}, error) {
	return func(ctx context.Context, input *CreateSubscribeInput) (*struct{}, error) {
		l := NewCreateSubscribeLogic(ctx, deps)
		if err := l.CreateSubscribe(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type CreateSubscribeLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// NewCreateSubscribeLogic Create subscribe
func NewCreateSubscribeLogic(ctx context.Context, deps Deps) *CreateSubscribeLogic {
	return &CreateSubscribeLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *CreateSubscribeLogic) CreateSubscribe(req *types.CreateSubscribeRequest) error {
	discount := ""
	if len(req.Discount) > 0 {
		val, _ := json.Marshal(req.Discount)
		discount = string(val)
	}
	sub := &subscribe.Subscribe{
		Id:                0,
		Name:              req.Name,
		Language:          req.Language,
		Description:       req.Description,
		UnitPrice:         req.UnitPrice,
		UnitTime:          req.UnitTime,
		Discount:          discount,
		Replacement:       req.Replacement,
		Inventory:         req.Inventory,
		Traffic:           req.Traffic,
		SpeedLimit:        req.SpeedLimit,
		DeviceLimit:       req.DeviceLimit,
		Quota:             req.Quota,
		Nodes:             tool.Int64SliceToString(req.Nodes),
		NodeTags:          tool.StringSliceToString(req.NodeTags),
		Show:              req.Show,
		Sell:              req.Sell,
		Sort:              0,
		DeductionRatio:    req.DeductionRatio,
		AllowDeduction:    req.AllowDeduction,
		ResetCycle:        req.ResetCycle,
		RenewalReset:      req.RenewalReset,
		ShowOriginalPrice: req.ShowOriginalPrice,
	}
	err := l.deps.SubscribeModel.Insert(l.ctx, sub)
	if err != nil {
		l.Error("[CreateSubscribeLogic] create subscribe error: ", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseInsertError), "create subscribe error: %v", err.Error())
	}

	return nil
}
