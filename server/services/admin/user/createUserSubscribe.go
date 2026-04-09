package user

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/perfect-panel/server/models/user"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/uuidx"
	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/pkg/errors"
	"time"
)

type CreateUserSubscribeInput struct {
	Body types.CreateUserSubscribeRequest
}

func CreateUserSubscribeHandler(deps Deps) func(context.Context, *CreateUserSubscribeInput) (*struct{}, error) {
	return func(ctx context.Context, input *CreateUserSubscribeInput) (*struct{}, error) {
		l := NewCreateUserSubscribeLogic(ctx, deps)
		if err := l.CreateUserSubscribe(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type CreateUserSubscribeLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Create user subcribe
func NewCreateUserSubscribeLogic(ctx context.Context, deps Deps) *CreateUserSubscribeLogic {
	return &CreateUserSubscribeLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *CreateUserSubscribeLogic) CreateUserSubscribe(req *types.CreateUserSubscribeRequest) error {
	// validate user
	userInfo, err := l.deps.UserModel.FindOne(l.ctx, req.UserId)
	if err != nil {
		l.Errorw("FindOne error", logger.Field("error", err.Error()), logger.Field("userId", req.UserId))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "FindOne error: %v", err.Error())
	}
	subs, err := l.deps.UserModel.QueryUserSubscribe(l.ctx, req.UserId)
	if err != nil {
		l.Errorw("QueryUserSubscribe error", logger.Field("error", err.Error()), logger.Field("userId", req.UserId))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "QueryUserSubscribe error: %v", err.Error())
	}
	if len(subs) >= 1 && l.deps.Config.Subscribe.SingleModel {
		return errors.Wrapf(xerr.NewErrCode(xerr.SingleSubscribeModeExceedsLimit), "Single subscribe mode exceeds limit")
	}
	sub, err := l.deps.SubscribeModel.FindOne(l.ctx, req.SubscribeId)
	if err != nil {
		l.Errorw("FindOne error", logger.Field("error", err.Error()), logger.Field("subscribeId", req.SubscribeId))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "FindOne error: %v", err.Error())
	}
	if req.Traffic == 0 {
		req.Traffic = sub.Traffic
	}

	userSub := user.Subscribe{
		UserId:      req.UserId,
		SubscribeId: req.SubscribeId,
		StartTime:   time.Now(),
		ExpireTime:  time.UnixMilli(req.ExpiredAt),
		Traffic:     req.Traffic,
		Download:    0,
		Upload:      0,
		Token:       uuidx.SubscribeToken(fmt.Sprintf("adminCreate:%d", time.Now().UnixMilli())),
		UUID:        uuid.New().String(),
		Status:      1,
	}
	if err = l.deps.UserModel.InsertSubscribe(l.ctx, &userSub); err != nil {
		l.Errorw("InsertSubscribe error", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseInsertError), "InsertSubscribe error: %v", err.Error())
	}

	err = l.deps.UserModel.UpdateUserCache(l.ctx, userInfo)
	if err != nil {
		l.Errorw("UpdateUserCache error", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseInsertError), "UpdateUserCache error: %v", err.Error())
	}

	err = l.deps.SubscribeModel.ClearCache(l.ctx, userSub.SubscribeId)
	if err != nil {
		logger.Errorw("ClearSubscribe error", logger.Field("error", err.Error()))
	}
	return nil
}
