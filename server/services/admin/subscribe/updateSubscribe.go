package subscribe

import (
	"context"
	"encoding/json"
	"github.com/perfect-panel/server/models/subscribe"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/modules/verify/device"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type UpdateSubscribeInput struct {
	Body types.UpdateSubscribeRequest
}

func UpdateSubscribeHandler(svcCtx *svc.ServiceContext) func(context.Context, *UpdateSubscribeInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateSubscribeInput) (*struct{}, error) {
		l := NewUpdateSubscribeLogic(ctx, svcCtx)
		if err := l.UpdateSubscribe(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type UpdateSubscribeLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Update subscribe
func NewUpdateSubscribeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSubscribeLogic {
	return &UpdateSubscribeLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateSubscribeLogic) UpdateSubscribe(req *types.UpdateSubscribeRequest) error {
	// Query the database to get the subscribe information
	_, err := l.svcCtx.SubscribeModel.FindOne(l.ctx, req.Id)
	if err != nil {
		l.Logger.Error("[UpdateSubscribe] Database query error", logger.Field("error", err.Error()), logger.Field("subscribe_id", req.Id))
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
	err = l.svcCtx.SubscribeModel.Update(l.ctx, sub)
	if err != nil {
		l.Logger.Error("[UpdateSubscribe] update subscribe failed", logger.Field("error", err.Error()), logger.Field("subscribe", sub))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseUpdateError), "update subscribe error: %v", err.Error())
	}
	l.svcCtx.DeviceManager.Broadcast(device.SubscribeUpdate)
	return nil
}
