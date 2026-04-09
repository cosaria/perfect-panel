package user

import (
	"context"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/pkg/errors"
	"os"
	"strings"
)

type DeleteUserInput struct {
	Body types.GetDetailRequest
}

func DeleteUserHandler(deps Deps) func(context.Context, *DeleteUserInput) (*struct{}, error) {
	return func(ctx context.Context, input *DeleteUserInput) (*struct{}, error) {
		l := NewDeleteUserLogic(ctx, deps)
		if err := l.DeleteUser(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type DeleteUserLogic struct {
	ctx  context.Context
	deps Deps
	logger.Logger
}

func NewDeleteUserLogic(ctx context.Context, deps Deps) *DeleteUserLogic {
	return &DeleteUserLogic{
		ctx:    ctx,
		deps:   deps,
		Logger: logger.WithContext(ctx),
	}
}

func (l *DeleteUserLogic) DeleteUser(req *types.GetDetailRequest) error {
	isDemo := strings.ToLower(os.Getenv("PPANEL_MODE")) == "demo"

	if req.Id == 2 && isDemo {
		return errors.Wrapf(xerr.NewErrCodeMsg(503, "Demo mode does not allow deletion of the admin user"), "delete user failed: cannot delete admin user in demo mode")
	}
	err := l.deps.UserModel.Delete(l.ctx, req.Id)
	if err != nil {
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseDeletedError), "delete user error: %v", err.Error())
	}
	return nil
}
