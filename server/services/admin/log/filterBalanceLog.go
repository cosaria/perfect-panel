package log

import (
	"context"
	"github.com/perfect-panel/server/models/log"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type FilterBalanceLogInput struct {
	types.FilterBalanceLogRequest
}

type FilterBalanceLogOutput struct {
	Body *types.FilterBalanceLogResponse
}

func FilterBalanceLogHandler(svcCtx *svc.ServiceContext) func(context.Context, *FilterBalanceLogInput) (*FilterBalanceLogOutput, error) {
	return func(ctx context.Context, input *FilterBalanceLogInput) (*FilterBalanceLogOutput, error) {
		l := NewFilterBalanceLogLogic(ctx, svcCtx)
		resp, err := l.FilterBalanceLog(&input.FilterBalanceLogRequest)
		if err != nil {
			return nil, err
		}
		return &FilterBalanceLogOutput{Body: resp}, nil
	}
}

type FilterBalanceLogLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewFilterBalanceLogLogic Filter balance log
func NewFilterBalanceLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FilterBalanceLogLogic {
	return &FilterBalanceLogLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FilterBalanceLogLogic) FilterBalanceLog(req *types.FilterBalanceLogRequest) (resp *types.FilterBalanceLogResponse, err error) {
	data, total, err := l.svcCtx.LogModel.FilterSystemLog(l.ctx, &log.FilterParams{
		Page:     req.Page,
		Size:     req.Size,
		Type:     log.TypeBalance.Uint8(),
		Data:     req.Date,
		ObjectID: req.UserId,
	})

	if err != nil {
		l.Errorw("[FilterBalanceLog] Query User Balance Log Error:", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "Query User Balance Log Error")
	}

	list := make([]types.BalanceLog, 0)
	for _, datum := range data {
		var content log.Balance
		if err = content.Unmarshal([]byte(datum.Content)); err != nil {
			l.Errorf("[QueryUserBalanceLog] unmarshal balance log content failed: %v", err.Error())
			continue
		}
		list = append(list, types.BalanceLog{
			UserId:    datum.ObjectID,
			Amount:    content.Amount,
			Type:      content.Type,
			OrderNo:   content.OrderNo,
			Balance:   content.Balance,
			Timestamp: content.Timestamp,
		})
	}

	return &types.FilterBalanceLogResponse{
		Total: total,
		List:  list,
	}, nil
}
