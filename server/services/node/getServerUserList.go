package server

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/perfect-panel/server/models/node"
	"github.com/perfect-panel/server/models/subscribe"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/modules/util/uuidx"
	"github.com/perfect-panel/server/routers/response"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"strings"
)

func GetServerUserListHandler(svcCtx *svc.ServiceContext) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req types.GetServerUserListRequest
		_ = c.ShouldBind(&req)
		_ = c.ShouldBindQuery(&req.ServerCommon)
		validateErr := svcCtx.Validate(&req)
		if validateErr != nil {
			response.ParamErrorResult(c, validateErr)
			return
		}

		l := NewGetServerUserListLogic(c, svcCtx)
		resp, err := l.GetServerUserList(&req)
		if err != nil {
			if errors.Is(err, xerr.StatusNotModified) {
				c.String(304, "Not Modified")
				return
			}
			c.String(404, "Not Found")
			return
		}
		c.JSON(200, resp)
	}
}

type GetServerUserListLogic struct {
	logger.Logger
	ctx    *gin.Context
	svcCtx *svc.ServiceContext
}

// NewGetServerUserListLogic Get user list
func NewGetServerUserListLogic(ctx *gin.Context, svcCtx *svc.ServiceContext) *GetServerUserListLogic {
	return &GetServerUserListLogic{
		Logger: logger.WithContext(ctx.Request.Context()),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetServerUserListLogic) GetServerUserList(req *types.GetServerUserListRequest) (resp *types.GetServerUserListResponse, err error) {
	cacheKey := fmt.Sprintf("%s%d", node.ServerUserListCacheKey, req.ServerId)
	cache, err := l.svcCtx.Redis.Get(l.ctx, cacheKey).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	if cache != "" {
		etag := tool.GenerateETag([]byte(cache))
		resp = &types.GetServerUserListResponse{}
		//  Check If-None-Match header
		if match := l.ctx.GetHeader("If-None-Match"); match == etag {
			return nil, xerr.StatusNotModified
		}
		l.ctx.Header("ETag", etag)
		err = json.Unmarshal([]byte(cache), resp)
		if err != nil {
			l.Errorw("[ServerUserListCacheKey] json unmarshal error", logger.Field("error", err.Error()))
			return nil, err
		}
		return resp, nil
	}
	server, err := l.svcCtx.NodeModel.FindOneServer(l.ctx, req.ServerId)
	if err != nil {
		return nil, err
	}

	_, nodes, err := l.svcCtx.NodeModel.FilterNodeList(l.ctx, &node.FilterNodeParams{
		Page:     1,
		Size:     1000,
		ServerId: []int64{server.Id},
		Protocol: req.Protocol,
	})
	if err != nil {
		l.Errorw("FilterNodeList error", logger.Field("error", err.Error()))
		return nil, err
	}
	var nodeTag []string
	var nodeIds []int64
	for _, n := range nodes {
		nodeIds = append(nodeIds, n.Id)
		if n.Tags != "" {
			nodeTag = append(nodeTag, strings.Split(n.Tags, ",")...)
		}
	}

	_, subs, err := l.svcCtx.SubscribeModel.FilterList(l.ctx, &subscribe.FilterParams{
		Page: 1,
		Size: 9999,
		Node: nodeIds,
		Tags: nodeTag,
	})
	if err != nil {
		l.Errorw("QuerySubscribeIdsByServerIdAndServerGroupId error", logger.Field("error", err.Error()))
		return nil, err
	}
	if len(subs) == 0 {
		return &types.GetServerUserListResponse{
			Users: []types.ServerUser{
				{
					Id:   1,
					UUID: uuidx.NewUUID().String(),
				},
			},
		}, nil
	}
	users := make([]types.ServerUser, 0)
	for _, sub := range subs {
		data, err := l.svcCtx.UserModel.FindUsersSubscribeBySubscribeId(l.ctx, sub.Id)
		if err != nil {
			return nil, err
		}
		for _, datum := range data {
			users = append(users, types.ServerUser{
				Id:          datum.Id,
				UUID:        datum.UUID,
				SpeedLimit:  sub.SpeedLimit,
				DeviceLimit: sub.DeviceLimit,
			})
		}
	}
	if len(users) == 0 {
		users = append(users, types.ServerUser{
			Id:   1,
			UUID: uuidx.NewUUID().String(),
		})
	}
	resp = &types.GetServerUserListResponse{
		Users: users,
	}
	val, _ := json.Marshal(resp)
	etag := tool.GenerateETag(val)
	l.ctx.Header("ETag", etag)
	err = l.svcCtx.Redis.Set(l.ctx, cacheKey, string(val), -1).Err()
	if err != nil {
		l.Errorw("[ServerUserListCacheKey] redis set error", logger.Field("error", err.Error()))
	}
	//  Check If-None-Match header
	if match := l.ctx.GetHeader("If-None-Match"); match == etag {
		return nil, xerr.StatusNotModified
	}
	return resp, nil
}
