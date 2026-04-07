package subscribe

import (
	"context"
	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/models/node"
	"github.com/perfect-panel/server/models/user"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
	"strings"
	"time"
)

type QueryUserSubscribeNodeListOutput struct {
	Body *types.QueryUserSubscribeNodeListResponse
}

func QueryUserSubscribeNodeListHandler(deps Deps) func(context.Context, *struct{}) (*QueryUserSubscribeNodeListOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*QueryUserSubscribeNodeListOutput, error) {
		l := NewQueryUserSubscribeNodeListLogic(ctx, deps)
		resp, err := l.QueryUserSubscribeNodeList()
		if err != nil {
			return nil, err
		}
		return &QueryUserSubscribeNodeListOutput{Body: resp}, nil
	}
}

type QueryUserSubscribeNodeListLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Get user subscribe node info
func NewQueryUserSubscribeNodeListLogic(ctx context.Context, deps Deps) *QueryUserSubscribeNodeListLogic {
	return &QueryUserSubscribeNodeListLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *QueryUserSubscribeNodeListLogic) QueryUserSubscribeNodeList() (resp *types.QueryUserSubscribeNodeListResponse, err error) {
	u, ok := l.ctx.Value(config.CtxKeyUser).(*user.User)
	if !ok {
		logger.Error("current user is not found in context")
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.InvalidAccess), "Invalid Access")
	}

	userSubscribes, err := l.deps.UserModel.QueryUserSubscribe(l.ctx, u.Id, 1, 2)
	if err != nil {
		logger.Errorw("failed to query user subscribe", logger.Field("error", err.Error()), logger.Field("user_id", u.Id))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "DB_ERROR")
	}

	resp = &types.QueryUserSubscribeNodeListResponse{}
	for _, us := range userSubscribes {
		userSubscribe, err := l.getUserSubscribe(us.Token)
		if err != nil {
			l.Errorw("[SubscribeLogic] Get user subscribe failed", logger.Field("error", err.Error()), logger.Field("token", userSubscribe.Token))
			return nil, err
		}
		nodes, err := l.getServers(userSubscribe)
		if err != nil {
			return nil, err
		}
		userSubscribeInfo := types.UserSubscribeInfo{
			Id:          userSubscribe.Id,
			Nodes:       nodes,
			Traffic:     userSubscribe.Traffic,
			Upload:      userSubscribe.Upload,
			Download:    userSubscribe.Download,
			Token:       userSubscribe.Token,
			UserId:      userSubscribe.UserId,
			OrderId:     userSubscribe.OrderId,
			SubscribeId: userSubscribe.SubscribeId,
			StartTime:   userSubscribe.StartTime.Unix(),
			ExpireTime:  userSubscribe.ExpireTime.Unix(),
			Status:      userSubscribe.Status,
			CreatedAt:   userSubscribe.CreatedAt.Unix(),
			UpdatedAt:   userSubscribe.UpdatedAt.Unix(),
		}

		if userSubscribe.FinishedAt != nil {
			userSubscribeInfo.FinishedAt = userSubscribe.FinishedAt.Unix()
		}

		cfg := l.deps.currentConfig()
		if cfg.Register.EnableTrial && cfg.Register.TrialSubscribe == userSubscribe.SubscribeId {
			userSubscribeInfo.IsTryOut = true
		}

		resp.List = append(resp.List, userSubscribeInfo)
	}

	return
}

func (l *QueryUserSubscribeNodeListLogic) getServers(userSub *user.Subscribe) (userSubscribeNodes []*types.UserSubscribeNodeInfo, err error) {
	userSubscribeNodes = make([]*types.UserSubscribeNodeInfo, 0)
	if l.isSubscriptionExpired(userSub) {
		return l.createExpiredServers(), nil
	}

	subDetails, err := l.deps.SubscribeModel.FindOne(l.ctx, userSub.SubscribeId)
	if err != nil {
		l.Errorw("[Generate Subscribe]find subscribe details error: %v", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "find subscribe details error: %v", err.Error())
	}
	nodeIds := tool.StringToInt64Slice(subDetails.Nodes)
	tags := strings.Split(subDetails.NodeTags, ",")

	l.Debugf("[Generate Subscribe]nodes: %v, NodeTags: %v", nodeIds, tags)

	enable := true

	_, nodes, err := l.deps.NodeModel.FilterNodeList(l.ctx, &node.FilterNodeParams{
		Page:    0,
		Size:    1000,
		NodeId:  nodeIds,
		Enabled: &enable, // Only get enabled nodes
	})

	if len(nodes) > 0 {
		var serverMapIds = make(map[int64]*node.Server)
		for _, n := range nodes {
			serverMapIds[n.ServerId] = nil
		}
		var serverIds []int64
		for k := range serverMapIds {
			serverIds = append(serverIds, k)
		}

		servers, err := l.deps.NodeModel.QueryServerList(l.ctx, serverIds)
		if err != nil {
			l.Errorw("[Generate Subscribe]find server details error: %v", logger.Field("error", err.Error()))
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "find server details error: %v", err.Error())
		}

		for _, s := range servers {
			serverMapIds[s.Id] = s
		}

		for _, n := range nodes {
			server := serverMapIds[n.ServerId]
			if server == nil {
				continue
			}
			userSubscribeNode := &types.UserSubscribeNodeInfo{
				Id:        n.Id,
				Name:      n.Name,
				Uuid:      userSub.UUID,
				Protocol:  n.Protocol,
				Port:      n.Port,
				Address:   n.Address,
				Tags:      strings.Split(n.Tags, ","),
				Country:   server.Country,
				City:      server.City,
				CreatedAt: n.CreatedAt.Unix(),
			}
			userSubscribeNodes = append(userSubscribeNodes, userSubscribeNode)
		}
	}

	l.Debugf("[Query Subscribe]found servers: %v", len(nodes))

	if err != nil {
		l.Errorw("[Generate Subscribe]find server details error: %v", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "find server details error: %v", err.Error())
	}
	logger.Debugf("[Generate Subscribe]found servers: %v", len(nodes))
	return userSubscribeNodes, nil
}

func (l *QueryUserSubscribeNodeListLogic) isSubscriptionExpired(userSub *user.Subscribe) bool {
	return userSub.ExpireTime.Unix() < time.Now().Unix() && userSub.ExpireTime.Unix() != 0
}

func (l *QueryUserSubscribeNodeListLogic) createExpiredServers() []*types.UserSubscribeNodeInfo {
	return nil
}

func (l *QueryUserSubscribeNodeListLogic) getUserSubscribe(token string) (*user.Subscribe, error) {
	userSub, err := l.deps.UserModel.FindOneSubscribeByToken(l.ctx, token)
	if err != nil {
		l.Infow("[Generate Subscribe]find subscribe error: %v", logger.Field("error", err.Error()), logger.Field("token", token))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "find subscribe error: %v", err.Error())
	}

	//  Ignore expiration check
	//if userSub.Status > 1 {
	//	l.Infow("[Generate Subscribe]subscribe is not available", logger.Field("status", int(userSub.Status)), logger.Field("token", token))
	//	return nil, errors.Wrapf(xerr.NewErrCode(xerr.SubscribeNotAvailable), "subscribe is not available")
	//}

	return userSub, nil
}
