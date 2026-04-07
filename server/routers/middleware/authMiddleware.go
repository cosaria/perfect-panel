package middleware

import (
	"context"
	"fmt"
	"strings"

	"github.com/perfect-panel/server/modules/infra/logger"

	"github.com/gin-gonic/gin"
	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/modules/auth/jwt"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/routers/response"
	"github.com/perfect-panel/server/svc"
	"github.com/pkg/errors"
)

func AuthMiddleware(svc *svc.ServiceContext) func(c *gin.Context) {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		jwtConfig := svc.Config.JwtAuth
		// get token from header
		token := c.GetHeader("Authorization")
		if token == "" {
			logger.WithContext(c.Request.Context()).Debug("[AuthMiddleware] Token Empty")
			response.HttpResult(c, nil, errors.Wrapf(xerr.NewErrCode(xerr.ErrorTokenEmpty), "Token Empty"))
			c.Abort()
			return
		}
		// parse token
		claims, err := jwt.ParseJwtToken(token, jwtConfig.AccessSecret)
		if err != nil {
			logger.WithContext(c.Request.Context()).Debug("[AuthMiddleware] ParseJwtToken", logger.Field("error", err.Error()), logger.Field("token", token))
			response.HttpResult(c, nil, errors.Wrapf(xerr.NewErrCode(xerr.ErrorTokenExpire), "Token Invalid"))
			c.Abort()
			return
		}

		loginType := ""
		if claims["LoginType"] != nil {
			loginType = claims["LoginType"].(string)
		}
		// get user id from token
		userId := int64(claims["UserId"].(float64))
		// get session id from token
		sessionId := claims["SessionId"].(string)
		// get session id from redis
		sessionIdCacheKey := fmt.Sprintf("%v:%v", config.SessionIdKey, sessionId)
		value, err := svc.Redis.Get(c, sessionIdCacheKey).Result()
		if err != nil {
			logger.WithContext(c.Request.Context()).Debug("[AuthMiddleware] Redis Get", logger.Field("error", err.Error()), logger.Field("sessionId", sessionId))
			response.HttpResult(c, nil, errors.Wrapf(xerr.NewErrCode(xerr.InvalidAccess), "Invalid Access"))
			c.Abort()
			return
		}

		//verify user id
		if value != fmt.Sprintf("%v", userId) {
			logger.WithContext(c.Request.Context()).Debug("[AuthMiddleware] Invalid Access", logger.Field("userId", userId), logger.Field("sessionId", sessionId))
			response.HttpResult(c, nil, errors.Wrapf(xerr.NewErrCode(xerr.InvalidAccess), "Invalid Access"))
			c.Abort()
			return
		}

		userInfo, err := svc.UserModel.FindOne(c, userId)
		if err != nil {
			logger.WithContext(c.Request.Context()).Debug("[AuthMiddleware] UserModel FindOne", logger.Field("error", err.Error()), logger.Field("userId", userId))
			response.HttpResult(c, nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "Database Query Error"))
			c.Abort()
			return
		}
		// admin verify
		paths := strings.Split(c.Request.URL.Path, "/")
		if tool.StringSliceContains(paths, "admin") && !*userInfo.IsAdmin {
			logger.WithContext(c.Request.Context()).Debug("[AuthMiddleware] Not Admin User", logger.Field("userId", userId), logger.Field("sessionId", sessionId))
			response.HttpResult(c, nil, errors.Wrapf(xerr.NewErrCode(xerr.InvalidAccess), "Invalid Access"))
			c.Abort()
			return
		}
		ctx = context.WithValue(ctx, config.LoginType, loginType)
		ctx = context.WithValue(ctx, config.CtxKeyUser, userInfo)
		ctx = context.WithValue(ctx, config.CtxKeySessionID, sessionId)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
