package log

import (
	"context"
	"github.com/perfect-panel/server/models/log"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type FilterCommissionLogInput struct {
	types.FilterCommissionLogRequest
}

type FilterCommissionLogOutput struct {
	Body *types.FilterCommissionLogResponse
}

func FilterCommissionLogHandler(deps Deps) func(context.Context, *FilterCommissionLogInput) (*FilterCommissionLogOutput, error) {
	return func(ctx context.Context, input *FilterCommissionLogInput) (*FilterCommissionLogOutput, error) {
		l := NewFilterCommissionLogLogic(ctx, deps)
		resp, err := l.FilterCommissionLog(&input.FilterCommissionLogRequest)
		if err != nil {
			return nil, err
		}
		return &FilterCommissionLogOutput{Body: resp}, nil
	}
}

type FilterCommissionLogLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// NewFilterCommissionLogLogic Filter commission log
func NewFilterCommissionLogLogic(ctx context.Context, deps Deps) *FilterCommissionLogLogic {
	return &FilterCommissionLogLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *FilterCommissionLogLogic) FilterCommissionLog(req *types.FilterCommissionLogRequest) (resp *types.FilterCommissionLogResponse, err error) {
	data, total, err := l.deps.LogModel.FilterSystemLog(l.ctx, &log.FilterParams{
		Page:     req.Page,
		Size:     req.Size,
		Data:     req.Date,
		Type:     log.TypeCommission.Uint8(),
		ObjectID: req.UserId,
	})
	if err != nil {
		l.Errorw("Query User Commission Log failed", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "Query User Commission Log failed")
	}
	var list []types.CommissionLog

	for _, datum := range data {
		var content log.Commission
		if err = content.Unmarshal([]byte(datum.Content)); err != nil {
			l.Errorf("unmarshal commission log content failed: %v", err.Error())
			continue
		}
		list = append(list, types.CommissionLog{
			UserId:    datum.ObjectID,
			Type:      content.Type,
			Amount:    content.Amount,
			OrderNo:   content.OrderNo,
			Timestamp: content.Timestamp,
		})
	}
	return &types.FilterCommissionLogResponse{
		Total: total,
		List:  list,
	}, nil
}
