package user

import (
	"context"
	"fmt"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/uuidx"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
	"time"
)

type ResetUserSubscribeTokenInput struct {
	Body types.ResetUserSubscribeTokenRequest
}

func ResetUserSubscribeTokenHandler(svcCtx *svc.ServiceContext) func(context.Context, *ResetUserSubscribeTokenInput) (*struct{}, error) {
	return func(ctx context.Context, input *ResetUserSubscribeTokenInput) (*struct{}, error) {
		l := NewResetUserSubscribeTokenLogic(ctx, svcCtx)
		if err := l.ResetUserSubscribeToken(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type ResetUserSubscribeTokenLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewResetUserSubscribeTokenLogic Reset user subscribe token
func NewResetUserSubscribeTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetUserSubscribeTokenLogic {
	return &ResetUserSubscribeTokenLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResetUserSubscribeTokenLogic) ResetUserSubscribeToken(req *types.ResetUserSubscribeTokenRequest) error {
	userSub, err := l.svcCtx.UserModel.FindOneSubscribe(l.ctx, req.UserSubscribeId)
	if err != nil {
		logger.Errorf("[ResetUserSubscribeToken] FindOneSubscribe error: %v", err.Error())
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "FindOneSubscribe error: %v", err.Error())
	}
	userSub.Token = uuidx.SubscribeToken(fmt.Sprintf("AdminUpdate:%d", time.Now().UnixMilli()))

	err = l.svcCtx.UserModel.UpdateSubscribe(l.ctx, userSub)
	if err != nil {
		logger.Errorf("[ResetUserSubscribeToken] UpdateSubscribe error: %v", err.Error())
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseUpdateError), "UpdateSubscribe error: %v", err.Error())
	}
	// Clear user subscribe cache
	if err = l.svcCtx.UserModel.ClearSubscribeCache(l.ctx, userSub); err != nil {
		l.Errorw("ClearSubscribeCache failed:", logger.Field("error", err.Error()), logger.Field("userSubscribeId", userSub.Id))
		return errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "ClearSubscribeCache failed: %v", err.Error())
	}
	// Clear subscribe cache
	if err = l.svcCtx.SubscribeModel.ClearCache(l.ctx, userSub.SubscribeId); err != nil {
		l.Errorw("failed to clear subscribe cache", logger.Field("error", err.Error()), logger.Field("subscribeId", userSub.SubscribeId))
		return errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "failed to clear subscribe cache: %v", err.Error())
	}
	return nil
}
