package user

import (
	"context"

	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/internal/platform/http/types"
)

type PreUnsubscribeInput struct {
	Body types.PreUnsubscribeRequest
}

type PreUnsubscribeOutput struct {
	Body *types.PreUnsubscribeResponse
}

func PreUnsubscribeHandler(deps Deps) func(context.Context, *PreUnsubscribeInput) (*PreUnsubscribeOutput, error) {
	return func(ctx context.Context, input *PreUnsubscribeInput) (*PreUnsubscribeOutput, error) {
		l := NewPreUnsubscribeLogic(ctx, deps)
		resp, err := l.PreUnsubscribe(&input.Body)
		if err != nil {
			return nil, err
		}
		return &PreUnsubscribeOutput{Body: resp}, nil
	}
}

type PreUnsubscribeLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// NewPreUnsubscribeLogic Pre Unsubscribe
func NewPreUnsubscribeLogic(ctx context.Context, deps Deps) *PreUnsubscribeLogic {
	return &PreUnsubscribeLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *PreUnsubscribeLogic) PreUnsubscribe(req *types.PreUnsubscribeRequest) (resp *types.PreUnsubscribeResponse, err error) {
	remainingAmount, err := CalculateRemainingAmount(l.ctx, l.deps, req.Id)
	if err != nil {
		l.Errorw("[PreUnsubscribeLogic] Calculate Remaining Amount Error:", logger.Field("err", err.Error()))
		return nil, err
	}
	return &types.PreUnsubscribeResponse{
		DeductionAmount: remainingAmount,
	}, nil
}
