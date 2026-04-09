package log

import (
	"context"

	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/internal/platform/persistence/log"
	"github.com/perfect-panel/server/internal/platform/support/logger"
	"github.com/perfect-panel/server/internal/platform/support/xerr"
	"github.com/pkg/errors"
)

type FilterGiftLogInput struct {
	types.FilterGiftLogRequest
}

type FilterGiftLogOutput struct {
	Body *types.FilterGiftLogResponse
}

func FilterGiftLogHandler(deps Deps) func(context.Context, *FilterGiftLogInput) (*FilterGiftLogOutput, error) {
	return func(ctx context.Context, input *FilterGiftLogInput) (*FilterGiftLogOutput, error) {
		l := NewFilterGiftLogLogic(ctx, deps)
		resp, err := l.FilterGiftLog(&input.FilterGiftLogRequest)
		if err != nil {
			return nil, err
		}
		return &FilterGiftLogOutput{Body: resp}, nil
	}
}

type FilterGiftLogLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Filter gift log
func NewFilterGiftLogLogic(ctx context.Context, deps Deps) *FilterGiftLogLogic {
	return &FilterGiftLogLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *FilterGiftLogLogic) FilterGiftLog(req *types.FilterGiftLogRequest) (resp *types.FilterGiftLogResponse, err error) {
	data, total, err := l.deps.LogModel.FilterSystemLog(l.ctx, &log.FilterParams{
		Page:     req.Page,
		Size:     req.Size,
		Type:     log.TypeGift.Uint8(),
		ObjectID: req.UserId,
		Data:     req.Date,
		Search:   req.Search,
	})

	if err != nil {
		l.Errorf("[FilterGiftLog] failed to filter system log: %v", err.Error())
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "failed to filter system log: %v", err.Error())
	}

	var list []types.GiftLog
	for _, datum := range data {
		var content log.Gift
		err = content.Unmarshal([]byte(datum.Content))
		if err != nil {
			l.Errorf("[FilterGiftLog] failed to unmarshal content: %v", err.Error())
			continue
		}
		list = append(list, types.GiftLog{
			Type:        content.Type,
			UserId:      datum.ObjectID,
			OrderNo:     content.OrderNo,
			SubscribeId: content.SubscribeId,
			Amount:      content.Amount,
			Balance:     content.Balance,
			Remark:      content.Remark,
			Timestamp:   content.Timestamp,
		})
	}

	return &types.FilterGiftLogResponse{
		Total: total,
		List:  list,
	}, nil
}
