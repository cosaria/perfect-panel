package user

import (
	"context"
	"github.com/perfect-panel/server/models/user"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/notify/phone"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/pkg/errors"
)

type GetUserListInput struct {
	Body types.GetUserListRequest
}

type GetUserListOutput struct {
	Body *types.GetUserListResponse
}

func GetUserListHandler(deps Deps) func(context.Context, *GetUserListInput) (*GetUserListOutput, error) {
	return func(ctx context.Context, input *GetUserListInput) (*GetUserListOutput, error) {
		l := NewGetUserListLogic(ctx, deps)
		resp, err := l.GetUserList(&input.Body)
		if err != nil {
			return nil, err
		}
		return &GetUserListOutput{Body: resp}, nil
	}
}

type GetUserListLogic struct {
	ctx  context.Context
	deps Deps
	logger.Logger
}

func NewGetUserListLogic(ctx context.Context, deps Deps) *GetUserListLogic {
	return &GetUserListLogic{
		ctx:    ctx,
		deps:   deps,
		Logger: logger.WithContext(ctx),
	}
}
func (l *GetUserListLogic) GetUserList(req *types.GetUserListRequest) (*types.GetUserListResponse, error) {
	list, total, err := l.deps.UserModel.QueryPageList(l.ctx, req.Page, req.Size, &user.UserFilterParams{
		UserId:          req.UserId,
		Search:          req.Search,
		Unscoped:        req.Unscoped,
		SubscribeId:     req.SubscribeId,
		UserSubscribeId: req.UserSubscribeId,
		Order:           "DESC",
	})
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "GetUserListLogic failed: %v", err.Error())
	}

	userRespList := make([]types.User, 0, len(list))

	for _, item := range list {
		var u types.User
		tool.DeepCopy(&u, item)

		// 处理 AuthMethods
		authMethods := make([]types.UserAuthMethod, len(u.AuthMethods)) // 直接创建目标 slice
		for i, method := range u.AuthMethods {
			tool.DeepCopy(&authMethods[i], method)
			if method.AuthType == "mobile" {
				authMethods[i].AuthIdentifier = phone.FormatToInternational(method.AuthIdentifier)
			}
		}
		u.AuthMethods = authMethods

		userRespList = append(userRespList, u)
	}

	return &types.GetUserListResponse{
		Total: total,
		List:  userRespList,
	}, nil
}
