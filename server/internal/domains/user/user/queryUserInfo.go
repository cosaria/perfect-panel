package user

import (
	"context"

	"sort"

	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/internal/platform/notify/phone"
	"github.com/perfect-panel/server/internal/platform/persistence/user"
	"github.com/perfect-panel/server/internal/platform/support/logger"
	"github.com/perfect-panel/server/internal/platform/support/tool"
	"github.com/perfect-panel/server/internal/platform/support/xerr"
	"github.com/pkg/errors"
)

type QueryUserInfoOutput struct {
	Body *types.User
}

func QueryUserInfoHandler(deps Deps) func(context.Context, *struct{}) (*QueryUserInfoOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*QueryUserInfoOutput, error) {
		l := NewQueryUserInfoLogic(ctx, deps)
		resp, err := l.QueryUserInfo()
		if err != nil {
			return nil, err
		}
		return &QueryUserInfoOutput{Body: resp}, nil
	}
}

type QueryUserInfoLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Query User Info
func NewQueryUserInfoLogic(ctx context.Context, deps Deps) *QueryUserInfoLogic {
	return &QueryUserInfoLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *QueryUserInfoLogic) QueryUserInfo() (resp *types.User, err error) {
	resp = &types.User{}
	u, ok := l.ctx.Value(config.CtxKeyUser).(*user.User)
	if !ok {
		logger.Error("current user is not found in context")
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.InvalidAccess), "Invalid Access")
	}
	tool.DeepCopy(resp, u)

	var userMethods []types.UserAuthMethod
	for _, method := range resp.AuthMethods {
		var item types.UserAuthMethod
		tool.DeepCopy(&item, method)

		switch method.AuthType {
		case "mobile":
			item.AuthIdentifier = phone.MaskPhoneNumber(method.AuthIdentifier)
		case "email":
		default:
			item.AuthIdentifier = maskOpenID(method.AuthIdentifier)
		}
		userMethods = append(userMethods, item)
	}

	// 按照指定顺序排序：email第一位，mobile第二位，其他按原顺序
	sort.Slice(userMethods, func(i, j int) bool {
		return getAuthTypePriority(userMethods[i].AuthType) < getAuthTypePriority(userMethods[j].AuthType)
	})

	resp.AuthMethods = userMethods
	return resp, nil
}

// getAuthTypePriority 获取认证类型的排序优先级
// email: 1 (第一位)
// mobile: 2 (第二位)
// 其他类型: 100+ (后续位置)
func getAuthTypePriority(authType string) int {
	switch authType {
	case "email":
		return 1
	case "mobile":
		return 2
	default:
		return 100
	}
}

// maskOpenID 脱敏 OpenID，只保留前 3 和后 3 位
func maskOpenID(openID string) string {
	length := len(openID)
	if length <= 6 {
		return "***" // 如果 ID 太短，直接返回 "***"
	}

	// 计算中间需要被替换的 `*` 数量
	maskLength := length - 6
	mask := make([]byte, maskLength)
	for i := range mask {
		mask[i] = '*'
	}

	// 组合脱敏后的 OpenID
	return openID[:3] + string(mask) + openID[length-3:]
}
