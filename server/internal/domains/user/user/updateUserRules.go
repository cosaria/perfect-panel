package user

import (
	"context"
	"encoding/json"

	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/models/user"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/pkg/errors"
)

type UpdateUserRulesInput struct {
	Body types.UpdateUserRulesRequest
}

func UpdateUserRulesHandler(deps Deps) func(context.Context, *UpdateUserRulesInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateUserRulesInput) (*struct{}, error) {
		l := NewUpdateUserRulesLogic(ctx, deps)
		if err := l.UpdateUserRules(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type UpdateUserRulesLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// NewUpdateUserRulesLogic Update User Rules
func NewUpdateUserRulesLogic(ctx context.Context, deps Deps) *UpdateUserRulesLogic {
	return &UpdateUserRulesLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *UpdateUserRulesLogic) UpdateUserRules(req *types.UpdateUserRulesRequest) error {
	u, ok := l.ctx.Value(config.CtxKeyUser).(*user.User)
	if !ok {
		logger.Error("current user is not found in context")
		return errors.Wrapf(xerr.NewErrCode(xerr.InvalidAccess), "Invalid Access")
	}
	if len(req.Rules) > 0 {
		bytes, err := json.Marshal(req.Rules)
		if err != nil {
			l.Errorf("UpdateUserRulesLogic json marshal rules error: %v", err)
			return errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "json marshal rules failed: %v", err.Error())
		}
		u.Rules = string(bytes)
		err = l.deps.UserModel.Update(l.ctx, u)
		if err != nil {
			l.Errorf("UpdateUserRulesLogic UpdateUserRules error: %v", err)
			return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseUpdateError), "update user rules failed: %v", err.Error())
		}
	}
	return nil
}
