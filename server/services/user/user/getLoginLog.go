package user

import (
	"context"

	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/models/log"
	"github.com/perfect-panel/server/models/user"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type GetLoginLogInput struct {
	types.GetLoginLogRequest
}

type GetLoginLogOutput struct {
	Body *types.GetLoginLogResponse
}

func GetLoginLogHandler(deps Deps) func(context.Context, *GetLoginLogInput) (*GetLoginLogOutput, error) {
	return func(ctx context.Context, input *GetLoginLogInput) (*GetLoginLogOutput, error) {
		l := NewGetLoginLogLogic(ctx, deps)
		resp, err := l.GetLoginLog(&input.GetLoginLogRequest)
		if err != nil {
			return nil, err
		}
		return &GetLoginLogOutput{Body: resp}, nil
	}
}

type GetLoginLogLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Get Login Log
func NewGetLoginLogLogic(ctx context.Context, deps Deps) *GetLoginLogLogic {
	return &GetLoginLogLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *GetLoginLogLogic) GetLoginLog(req *types.GetLoginLogRequest) (resp *types.GetLoginLogResponse, err error) {
	u, ok := l.ctx.Value(config.CtxKeyUser).(*user.User)
	if !ok {
		logger.Error("current user is not found in context")
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.InvalidAccess), "Invalid Access")
	}
	data, total, err := l.deps.LogModel.FilterSystemLog(l.ctx, &log.FilterParams{
		Page:     req.Page,
		Size:     req.Size,
		Type:     log.TypeLogin.Uint8(),
		ObjectID: u.Id,
	})
	if err != nil {
		l.Errorw("find login log failed:", logger.Field("error", err.Error()), logger.Field("user_id", u.Id))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "find login log failed: %v", err.Error())
	}
	list := make([]types.UserLoginLog, 0)

	for _, datum := range data {
		var content log.Login
		if err = content.Unmarshal([]byte(datum.Content)); err != nil {
			l.Errorf("[GetUserLoginLogs] unmarshal login log content failed: %v", err.Error())
			continue
		}
		list = append(list, types.UserLoginLog{
			Id:        datum.Id,
			UserId:    datum.ObjectID,
			LoginIP:   content.LoginIP,
			UserAgent: content.UserAgent,
			Success:   content.Success,
			Timestamp: datum.CreatedAt.UnixMilli(),
		})
	}

	return &types.GetLoginLogResponse{
		Total: total,
		List:  list,
	}, nil
}
