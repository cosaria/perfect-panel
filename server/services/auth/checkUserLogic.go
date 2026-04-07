package auth

import (
	"context"

	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type CheckUserLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewCheckUserLogic Check user is exist
func NewCheckUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckUserLogic {
	return &CheckUserLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CheckUserLogic) CheckUser(req *types.CheckUserRequest) (resp *types.CheckUserResponse, err error) {
	authMethod, err := l.svcCtx.UserModel.FindUserAuthMethodByOpenID(l.ctx, "email", req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "find user by email error: %v", err.Error())
	}
	return &types.CheckUserResponse{
		Exist: authMethod.UserId != 0,
	}, nil
}
