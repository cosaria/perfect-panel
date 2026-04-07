package middleware

import (
	"context"
	"fmt"
	"net/http"
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
			response.WriteProblem(c, response.NewPublicProblem(http.StatusUnauthorized, response.ProblemTypeUnauthorized, http.StatusText(http.StatusUnauthorized)))
			c.Abort()
			return
		}
		// parse token
		claims, err := jwt.ParseJwtToken(token, jwtConfig.AccessSecret)
		if err != nil {
			logger.WithContext(c.Request.Context()).Debug("[AuthMiddleware] ParseJwtToken", logger.Field("error", err.Error()), logger.Field("token", token))
			response.WriteProblem(c, response.NewPublicProblem(http.StatusUnauthorized, response.ProblemTypeUnauthorized, http.StatusText(http.StatusUnauthorized)))
			c.Abort()
			return
		}

		loginType := ""
		if claims["LoginType"] != nil {
			var ok bool
			loginType, ok = claims["LoginType"].(string)
			if !ok {
				logger.WithContext(c.Request.Context()).Debug("[AuthMiddleware] Invalid LoginType claim", logger.Field("claim", claims["LoginType"]))
				response.WriteProblem(c, response.NewPublicProblem(http.StatusUnauthorized, response.ProblemTypeUnauthorized, http.StatusText(http.StatusUnauthorized)))
				c.Abort()
				return
			}
		}
		rawUserId, ok := claims["UserId"].(float64)
		if !ok {
			logger.WithContext(c.Request.Context()).Debug("[AuthMiddleware] Invalid UserId claim", logger.Field("claim", claims["UserId"]))
			response.WriteProblem(c, response.NewPublicProblem(http.StatusUnauthorized, response.ProblemTypeUnauthorized, http.StatusText(http.StatusUnauthorized)))
			c.Abort()
			return
		}
		userId := int64(rawUserId)
		sessionId, ok := claims["SessionId"].(string)
		if !ok || sessionId == "" {
			logger.WithContext(c.Request.Context()).Debug("[AuthMiddleware] Invalid SessionId claim", logger.Field("claim", claims["SessionId"]))
			response.WriteProblem(c, response.NewPublicProblem(http.StatusUnauthorized, response.ProblemTypeUnauthorized, http.StatusText(http.StatusUnauthorized)))
			c.Abort()
			return
		}
		// get session id from redis
		sessionIdCacheKey := fmt.Sprintf("%v:%v", config.SessionIdKey, sessionId)
		value, err := svc.Redis.Get(c, sessionIdCacheKey).Result()
		if err != nil {
			logger.WithContext(c.Request.Context()).Debug("[AuthMiddleware] Redis Get", logger.Field("error", err.Error()), logger.Field("sessionId", sessionId))
			response.WriteProblem(c, response.NewPublicProblem(http.StatusUnauthorized, response.ProblemTypeUnauthorized, http.StatusText(http.StatusUnauthorized)))
			c.Abort()
			return
		}

		//verify user id
		if value != fmt.Sprintf("%v", userId) {
			logger.WithContext(c.Request.Context()).Debug("[AuthMiddleware] Invalid Access", logger.Field("userId", userId), logger.Field("sessionId", sessionId))
			response.WriteProblem(c, response.NewPublicProblem(http.StatusUnauthorized, response.ProblemTypeUnauthorized, http.StatusText(http.StatusUnauthorized)))
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
			response.WriteProblem(c, response.NewPublicProblem(http.StatusForbidden, response.ProblemTypeForbidden, http.StatusText(http.StatusForbidden)))
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
