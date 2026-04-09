package user

import (
	"context"
	"os"
	"strings"

	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/internal/platform/support/logger"
	"github.com/perfect-panel/server/internal/platform/support/tool"
	"github.com/perfect-panel/server/internal/platform/support/xerr"
	"github.com/pkg/errors"
)

type BatchDeleteUserInput struct {
	Body types.BatchDeleteUserRequest
}

func BatchDeleteUserHandler(deps Deps) func(context.Context, *BatchDeleteUserInput) (*struct{}, error) {
	return func(ctx context.Context, input *BatchDeleteUserInput) (*struct{}, error) {
		l := NewBatchDeleteUserLogic(ctx, deps)
		if err := l.BatchDeleteUser(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type BatchDeleteUserLogic struct {
	ctx  context.Context
	deps Deps
	logger.Logger
}

func NewBatchDeleteUserLogic(ctx context.Context, deps Deps) *BatchDeleteUserLogic {
	return &BatchDeleteUserLogic{
		ctx:    ctx,
		deps:   deps,
		Logger: logger.WithContext(ctx),
	}
}

func (l *BatchDeleteUserLogic) BatchDeleteUser(req *types.BatchDeleteUserRequest) error {
	isDemo := strings.ToLower(os.Getenv("PPANEL_MODE")) == "demo"

	if tool.Contains(req.Ids, 2) && isDemo {
		return errors.Wrapf(xerr.NewErrCodeMsg(503, "Demo mode does not allow deletion of the admin user"), "BatchDeleteUser failed: cannot delete admin user in demo mode")
	}

	err := l.deps.UserModel.BatchDeleteUser(l.ctx, req.Ids)
	if err != nil {
		l.Error("[BatchDeleteUserLogic] BatchDeleteUser failed: ", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseDeletedError), "BatchDeleteUser failed: %v", err.Error())
	}
	return nil
}
