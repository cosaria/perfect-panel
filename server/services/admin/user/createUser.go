package user

import (
	"context"
	"fmt"
	"github.com/perfect-panel/server/models/user"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/modules/util/uuidx"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

type CreateUserInput struct {
	Body types.CreateUserRequest
}

func CreateUserHandler(svcCtx *svc.ServiceContext) func(context.Context, *CreateUserInput) (*struct{}, error) {
	return func(ctx context.Context, input *CreateUserInput) (*struct{}, error) {
		l := NewCreateUserLogic(ctx, svcCtx)
		if err := l.CreateUser(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type CreateUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger.Logger
}

func NewCreateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateUserLogic {
	return &CreateUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logger.WithContext(ctx),
	}
}
func (l *CreateUserLogic) CreateUser(req *types.CreateUserRequest) error {
	if req.ReferCode == "" {
		// timestamp replaces user id
		req.ReferCode = uuidx.UserInviteCode(time.Now().UnixMicro())
	}
	if req.Password == "" {
		req.Password = req.Email
	}
	pwd := tool.EncodePassWord(req.Password)
	newUser := &user.User{
		Password:           pwd,
		Algo:               "default",
		ReferralPercentage: req.ReferralPercentage,
		OnlyFirstPurchase:  &req.OnlyFirstPurchase,
		ReferCode:          req.ReferCode,
		Balance:            req.Balance,
		IsAdmin:            &req.IsAdmin,
	}
	var ams []user.AuthMethods

	if req.TelephoneAreaCode != "" && req.Telephone != "" {
		phone := fmt.Sprintf("%s-%s", req.TelephoneAreaCode, req.Telephone)
		_, err := l.svcCtx.UserModel.FindUserAuthMethodByOpenID(l.ctx, "mobile", phone)
		if err == nil {
			return errors.Wrapf(xerr.NewErrCode(xerr.TelephoneExist), "telephone exist")
		}
		ams = append(ams, user.AuthMethods{
			AuthType:       "mobile",
			AuthIdentifier: phone,
		})
	}
	if req.Email != "" {
		_, err := l.svcCtx.UserModel.FindUserAuthMethodByOpenID(l.ctx, "email", req.Email)
		if err == nil {
			return errors.Wrapf(xerr.NewErrCode(xerr.EmailExist), "email exist")
		}
		ams = append(ams, user.AuthMethods{
			AuthType:       "email",
			AuthIdentifier: req.Email,
		})
	}

	newUser.AuthMethods = ams

	// todo: get product id and duration
	if req.RefererUser != "" {
		// get referer user id
		u, err := l.svcCtx.UserModel.FindOneByEmail(l.ctx, req.RefererUser)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.Wrapf(xerr.NewErrCode(xerr.UserNotExist), "referer user not found: %v", err.Error())
			}
			return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "find referer user failed: %v", err.Error())
		}
		newUser.RefererId = u.Id
	}

	err := l.svcCtx.UserModel.Insert(l.ctx, newUser)
	if err != nil {
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseInsertError), "insert user failed: %v", err.Error())
	}
	return nil
}
