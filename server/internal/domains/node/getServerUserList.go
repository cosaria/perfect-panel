package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/perfect-panel/server/internal/platform/http/response"
	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/internal/platform/persistence/node"
	"github.com/perfect-panel/server/internal/platform/support/logger"
	"github.com/perfect-panel/server/internal/platform/support/tool"
	"github.com/perfect-panel/server/internal/platform/support/xerr"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

func GetServerUserListHandler(deps Deps) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req types.GetServerUserListRequest
		_ = c.ShouldBind(&req)
		_ = c.ShouldBindQuery(&req.ServerCommon)
		validateErr := validateRequest(&req)
		if validateErr != nil {
			response.ParamErrorResult(c, validateErr)
			return
		}

		l := NewGetServerUserListLogic(c, deps)
		resp, err := l.GetServerUserList(&req)
		if err != nil {
			if errors.Is(err, xerr.ErrStatusNotModified) {
				c.Status(http.StatusNotModified)
				return
			}
			response.WriteProblem(c, response.NewPublicProblem(http.StatusBadGateway, response.ProblemTypeNodeUnavailable, "Node resource unavailable"))
			return
		}
		c.JSON(200, resp)
	}
}

type GetServerUserListLogic struct {
	logger.Logger
	ctx  *gin.Context
	deps Deps
}

// NewGetServerUserListLogic Get user list
func NewGetServerUserListLogic(ctx *gin.Context, deps Deps) *GetServerUserListLogic {
	return &GetServerUserListLogic{
		Logger: logger.WithContext(ctx.Request.Context()),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *GetServerUserListLogic) GetServerUserList(req *types.GetServerUserListRequest) (resp *types.GetServerUserListResponse, err error) {
	cacheKey := fmt.Sprintf("%s%d", node.ServerUserListCacheKey, req.ServerId)
	cache, err := l.deps.Redis.Get(l.ctx, cacheKey).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	if cache != "" {
		etag := tool.GenerateETag([]byte(cache))
		resp = &types.GetServerUserListResponse{}
		//  Check If-None-Match header
		if match := l.ctx.GetHeader("If-None-Match"); match == etag {
			return nil, xerr.ErrStatusNotModified
		}
		l.ctx.Header("ETag", etag)
		err = json.Unmarshal([]byte(cache), resp)
		if err != nil {
			l.Errorw("[ServerUserListCacheKey] json unmarshal error", logger.Field("error", err.Error()))
			return nil, err
		}
		return resp, nil
	}
	server, err := l.deps.NodeModel.FindOneServer(l.ctx, req.ServerId)
	if err != nil {
		return nil, err
	}

	users := make([]types.ServerUser, 0)
	assignedIds, err := l.deps.NodeModel.ListAssignedUserSubscriptionIDs(l.ctx, server.Id, req.Protocol)
	if err != nil {
		l.Errorw("ListAssignedUserSubscriptionIDs error", logger.Field("error", err.Error()))
		return nil, err
	}
	subscriptions, err := l.deps.UserModel.ActivateAndFindUserSubscribeDetailsByIDs(l.ctx, assignedIds)
	if err != nil {
		l.Errorw("ActivateAndFindUserSubscribeDetailsByIDs error", logger.Field("error", err.Error()))
		return nil, err
	}
	for _, data := range subscriptions {
		if data == nil {
			continue
		}
		speedLimit := int64(0)
		deviceLimit := int64(0)
		if data.Subscribe != nil {
			speedLimit = data.Subscribe.SpeedLimit
			deviceLimit = data.Subscribe.DeviceLimit
		}
		users = append(users, types.ServerUser{
			Id:          data.Id,
			UUID:        data.UUID,
			SpeedLimit:  speedLimit,
			DeviceLimit: deviceLimit,
		})
	}
	resp = &types.GetServerUserListResponse{
		Users: users,
	}
	val, _ := json.Marshal(resp)
	etag := tool.GenerateETag(val)
	l.ctx.Header("ETag", etag)
	err = l.deps.Redis.Set(l.ctx, cacheKey, string(val), -1).Err()
	if err != nil {
		l.Errorw("[ServerUserListCacheKey] redis set error", logger.Field("error", err.Error()))
	}
	//  Check If-None-Match header
	if match := l.ctx.GetHeader("If-None-Match"); match == etag {
		return nil, xerr.ErrStatusNotModified
	}
	return resp, nil
}
