package user

import (
	"context"

	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/internal/platform/persistence/user"
	"github.com/perfect-panel/server/internal/platform/support/logger"
	"github.com/perfect-panel/server/internal/platform/support/xerr"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type CreateUserAuthMethodInput struct {
	Body types.CreateUserAuthMethodRequest
}

func CreateUserAuthMethodHandler(deps Deps) func(context.Context, *CreateUserAuthMethodInput) (*struct{}, error) {
	return func(ctx context.Context, input *CreateUserAuthMethodInput) (*struct{}, error) {
		l := NewCreateUserAuthMethodLogic(ctx, deps)
		if err := l.CreateUserAuthMethod(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type CreateUserAuthMethodLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Create user auth method
func NewCreateUserAuthMethodLogic(ctx context.Context, deps Deps) *CreateUserAuthMethodLogic {
	return &CreateUserAuthMethodLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *CreateUserAuthMethodLogic) CreateUserAuthMethod(req *types.CreateUserAuthMethodRequest) error {
	err := l.deps.UserModel.Transaction(l.ctx, func(db *gorm.DB) error {
		var data *user.AuthMethods
		if err := db.Model(&user.AuthMethods{}).Where("`user_id` = ? AND `auth_type` = ?", req.UserId, req.AuthType).First(&data).Error; err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		data.UserId = req.UserId
		data.AuthType = req.AuthType
		data.AuthIdentifier = req.AuthIdentifier
		if err := db.Model(&user.AuthMethods{}).Save(&data).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		l.Errorw("[CreateUserAuthMethodLogic] Create User Auth Method Error:", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseInsertError), "Create User Auth Method Error")
	}
	return nil
}
