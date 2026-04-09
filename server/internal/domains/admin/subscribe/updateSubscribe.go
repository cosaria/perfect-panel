package subscribe

import (
	"context"
	"encoding/json"

	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/internal/platform/persistence/subscribe"
	"github.com/perfect-panel/server/internal/platform/support/logger"
	"github.com/perfect-panel/server/internal/platform/support/tool"
	"github.com/perfect-panel/server/internal/platform/support/verify/device"
	"github.com/perfect-panel/server/internal/platform/support/xerr"
	"github.com/pkg/errors"
)

type UpdateSubscribeInput struct {
	Body types.UpdateSubscribeRequest
}

func UpdateSubscribeHandler(deps Deps) func(context.Context, *UpdateSubscribeInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateSubscribeInput) (*struct{}, error) {
		l := NewUpdateSubscribeLogic(ctx, deps)
		if err := l.UpdateSubscribe(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type UpdateSubscribeLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Update subscribe
func NewUpdateSubscribeLogic(ctx context.Context, deps Deps) *UpdateSubscribeLogic {
	return &UpdateSubscribeLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *UpdateSubscribeLogic) UpdateSubscribe(req *types.UpdateSubscribeRequest) error {
	// Query the database to get the subscribe information
	_, err := l.deps.SubscribeModel.FindOne(l.ctx, req.Id)
	if err != nil {
		l.Error("[UpdateSubscribe] Database query error", logger.Field("error", err.Error()), logger.Field("subscribe_id", req.Id))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "get subscribe error: %v", err.Error())
	}
	discount := ""
	if len(req.Discount) > 0 {
		val, _ := json.Marshal(req.Discount)
		discount = string(val)
	}
	sub := &subscribe.Subscribe{
		Id:                req.Id,
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
		Sort:              req.Sort,
		DeductionRatio:    req.DeductionRatio,
		AllowDeduction:    req.AllowDeduction,
		ResetCycle:        req.ResetCycle,
		RenewalReset:      req.RenewalReset,
		ShowOriginalPrice: req.ShowOriginalPrice,
	}
	err = l.deps.SubscribeModel.Update(l.ctx, sub)
	if err != nil {
		l.Error("[UpdateSubscribe] update subscribe failed", logger.Field("error", err.Error()), logger.Field("subscribe", sub))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseUpdateError), "update subscribe error: %v", err.Error())
	}
	if l.deps.DeviceManager != nil {
		l.deps.DeviceManager.Broadcast(device.SubscribeUpdate)
	}
	return nil
}
