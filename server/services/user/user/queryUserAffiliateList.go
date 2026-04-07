package user

import (
	"context"

	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/models/user"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type QueryUserAffiliateListInput struct {
	types.QueryUserAffiliateListRequest
}

type QueryUserAffiliateListOutput struct {
	Body *types.QueryUserAffiliateListResponse
}

func QueryUserAffiliateListHandler(deps Deps) func(context.Context, *QueryUserAffiliateListInput) (*QueryUserAffiliateListOutput, error) {
	return func(ctx context.Context, input *QueryUserAffiliateListInput) (*QueryUserAffiliateListOutput, error) {
		l := NewQueryUserAffiliateListLogic(ctx, deps)
		resp, err := l.QueryUserAffiliateList(&input.QueryUserAffiliateListRequest)
		if err != nil {
			return nil, err
		}
		return &QueryUserAffiliateListOutput{Body: resp}, nil
	}
}

type QueryUserAffiliateListLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Query User Affiliate List
func NewQueryUserAffiliateListLogic(ctx context.Context, deps Deps) *QueryUserAffiliateListLogic {
	return &QueryUserAffiliateListLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *QueryUserAffiliateListLogic) QueryUserAffiliateList(req *types.QueryUserAffiliateListRequest) (resp *types.QueryUserAffiliateListResponse, err error) {
	u, ok := l.ctx.Value(config.CtxKeyUser).(*user.User)
	if !ok {
		logger.Error("current user is not found in context")
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.InvalidAccess), "Invalid Access")
	}
	var data []*user.User
	var total int64
	err = l.deps.UserModel.Transaction(l.ctx, func(db *gorm.DB) error {
		return db.Model(&user.User{}).Order("id desc").Where("referer_id = ?", u.Id).Count(&total).Limit(req.Size).Offset((req.Page - 1) * req.Size).Find(&data).Error
	})
	if err != nil {
		l.Errorw("Query User Affiliate List failed: %v", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "Query User Affiliate List failed: %v", err.Error())
	}

	list := make([]types.UserAffiliate, 0)
	for _, item := range data {
		list = append(list, types.UserAffiliate{
			Identifier:   GetAuthMethod(l, item).AuthIdentifier,
			Avatar:       item.Avatar,
			RegisteredAt: item.CreatedAt.UnixMilli(),
			Enable:       *item.Enable,
		})
	}
	return &types.QueryUserAffiliateListResponse{
		Total: total,
		List:  list,
	}, nil
}

func GetAuthMethod(l *QueryUserAffiliateListLogic, item *user.User) user.AuthMethods {
	authMethod := user.AuthMethods{}
	authMethods, errs := l.deps.UserModel.FindUserAuthMethods(l.ctx, item.Id)
	if errs == nil && len(authMethods) > 0 {
		for _, am := range authMethods {
			if am.AuthType == "6" || am.AuthType == "7" {
				authMethod = *am
				break
			}
		}
		if authMethod.AuthIdentifier == "" {
			authMethod = *authMethods[0]
		}

		hideTextLength := len(authMethod.AuthIdentifier) / 3
		if hideTextLength > 0 {
			authMethod.AuthIdentifier = authMethod.AuthIdentifier[0:hideTextLength] + "***" + authMethod.AuthIdentifier[hideTextLength*2:]
		}
	}
	return authMethod
}
