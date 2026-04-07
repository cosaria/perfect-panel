package user

import (
	"context"

	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type GetUserSubscribeByIdLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Get user subcribe by id
func NewGetUserSubscribeByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserSubscribeByIdLogic {
	return &GetUserSubscribeByIdLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserSubscribeByIdLogic) GetUserSubscribeById(req *types.GetUserSubscribeByIdRequest) (resp *types.UserSubscribeDetail, err error) {
	sub, err := l.svcCtx.UserModel.FindOneSubscribeDetailsById(l.ctx, req.Id)
	if err != nil {
		l.Errorw("[GetUserSubscribeByIdLogic] FindOneSubscribeDetailsById error", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "FindOneSubscribeDetailsById error: %v", err.Error())
	}
	var subscribeDetails types.UserSubscribeDetail
	tool.DeepCopy(&subscribeDetails, sub)
	return &subscribeDetails, nil
}
