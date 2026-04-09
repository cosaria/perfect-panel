package user

import (
	"context"

	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/internal/platform/support/logger"
	"github.com/perfect-panel/server/internal/platform/support/tool"
	"github.com/perfect-panel/server/internal/platform/support/xerr"
	"github.com/pkg/errors"
)

type GetUserSubscribeByIdInput struct {
	types.GetUserSubscribeByIdRequest
}

type GetUserSubscribeByIdOutput struct {
	Body *types.UserSubscribeDetail
}

func GetUserSubscribeByIdHandler(deps Deps) func(context.Context, *GetUserSubscribeByIdInput) (*GetUserSubscribeByIdOutput, error) {
	return func(ctx context.Context, input *GetUserSubscribeByIdInput) (*GetUserSubscribeByIdOutput, error) {
		l := NewGetUserSubscribeByIdLogic(ctx, deps)
		resp, err := l.GetUserSubscribeById(&input.GetUserSubscribeByIdRequest)
		if err != nil {
			return nil, err
		}
		return &GetUserSubscribeByIdOutput{Body: resp}, nil
	}
}

type GetUserSubscribeByIdLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Get user subcribe by id
func NewGetUserSubscribeByIdLogic(ctx context.Context, deps Deps) *GetUserSubscribeByIdLogic {
	return &GetUserSubscribeByIdLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *GetUserSubscribeByIdLogic) GetUserSubscribeById(req *types.GetUserSubscribeByIdRequest) (resp *types.UserSubscribeDetail, err error) {
	sub, err := l.deps.UserModel.FindOneSubscribeDetailsById(l.ctx, req.Id)
	if err != nil {
		l.Errorw("[GetUserSubscribeByIdLogic] FindOneSubscribeDetailsById error", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "FindOneSubscribeDetailsById error: %v", err.Error())
	}
	var subscribeDetails types.UserSubscribeDetail
	tool.DeepCopy(&subscribeDetails, sub)
	return &subscribeDetails, nil
}
