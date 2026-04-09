package server

import (
	"context"
	"time"

	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/internal/platform/persistence/node"
	"github.com/perfect-panel/server/internal/platform/support/logger"
	"github.com/perfect-panel/server/internal/platform/support/tool"
	"github.com/perfect-panel/server/internal/platform/support/xerr"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type FilterServerListInput struct {
	types.FilterServerListRequest
}

type FilterServerListOutput struct {
	Body *types.FilterServerListResponse
}

func FilterServerListHandler(deps Deps) func(context.Context, *FilterServerListInput) (*FilterServerListOutput, error) {
	return func(ctx context.Context, input *FilterServerListInput) (*FilterServerListOutput, error) {
		l := NewFilterServerListLogic(ctx, deps)
		resp, err := l.FilterServerList(&input.FilterServerListRequest)
		if err != nil {
			return nil, err
		}
		return &FilterServerListOutput{Body: resp}, nil
	}
}

type FilterServerListLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// NewFilterServerListLogic Filter Server List
func NewFilterServerListLogic(ctx context.Context, deps Deps) *FilterServerListLogic {
	return &FilterServerListLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *FilterServerListLogic) FilterServerList(req *types.FilterServerListRequest) (resp *types.FilterServerListResponse, err error) {
	total, data, err := l.deps.NodeModel.FilterServerList(l.ctx, &node.FilterParams{
		Page:   req.Page,
		Size:   req.Size,
		Search: req.Search,
	})
	if err != nil {
		l.Errorw("[FilterServerList] Query Database Error: ", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "[FilterServerList] Query Database Error")
	}

	list := make([]types.Server, 0)

	for _, datum := range data {
		var server types.Server
		tool.DeepCopy(&server, datum)

		// handler protocols
		var protocols []types.Protocol
		dst, err := datum.UnmarshalProtocols()
		if err != nil {
			l.Errorf("[FilterServerList] UnmarshalProtocols Error: %s", err.Error())
			continue
		}
		tool.DeepCopy(&protocols, dst)
		server.Protocols = protocols

		nodeStatus, err := l.deps.NodeModel.StatusCache(l.ctx, datum.Id)
		if err != nil {
			if !errors.Is(err, redis.Nil) {
				l.Errorw("[handlerServerStatus] GetNodeStatus Error: ", logger.Field("error", err.Error()), logger.Field("node_id", datum.Id))
			}
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "GetNodeStatus Error")
		}
		server.Status = types.ServerStatus{
			Mem:    nodeStatus.Mem,
			Cpu:    nodeStatus.Cpu,
			Disk:   nodeStatus.Disk,
			Online: l.handlerServerStatus(datum.Id, protocols),
			Status: l.handlerServerStaus(datum.LastReportedAt),
		}
		list = append(list, server)
	}

	return &types.FilterServerListResponse{
		List:  list,
		Total: total,
	}, nil
}

func (l *FilterServerListLogic) handlerServerStatus(id int64, protocols []types.Protocol) []types.ServerOnlineUser {
	result := make([]types.ServerOnlineUser, 0)

	for _, protocol := range protocols {
		// query online user
		data, err := l.deps.NodeModel.OnlineUserSubscribe(l.ctx, id, protocol.Type)
		if err != nil {
			if !errors.Is(err, redis.Nil) {
				l.Errorw("[handlerServerStatus] OnlineUserSubscribe Error: ", logger.Field("error", err.Error()), logger.Field("node_id", id), logger.Field("protocol", protocol.Type))
			}
			continue
		}
		if len(data) > 0 {
			for sub, online := range data {
				var ips []types.ServerOnlineIP
				for _, ip := range online {
					ips = append(ips, types.ServerOnlineIP{
						IP:       ip,
						Protocol: protocol.Type,
					})
				}

				result = append(result, types.ServerOnlineUser{
					IP:          ips,
					SubscribeId: sub,
				})
			}
		}
	}
	// merge same subscribe
	var mapResult = make(map[int64]types.ServerOnlineUser)
	for _, item := range result {
		if exist, ok := mapResult[item.SubscribeId]; ok {
			// merge
			exist.Traffic += item.Traffic
			exist.IP = append(exist.IP, item.IP...)
			mapResult[item.SubscribeId] = exist
		} else {
			// get subscribe info
			info, err := l.deps.UserModel.FindOneUserSubscribe(l.ctx, item.SubscribeId)
			if err != nil {
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					l.Errorw("[handlerServerStatus] FindOneSubscribe Error: ", logger.Field("error", err.Error()), logger.Field("subscribe_id", item.SubscribeId))
				}
				continue
			}
			data := types.ServerOnlineUser{
				IP:          item.IP,
				UserId:      info.UserId,
				Subscribe:   "",
				SubscribeId: item.SubscribeId,
				Traffic:     info.Download + info.Upload,
				ExpiredAt:   info.ExpireTime.UnixMilli(),
			}
			if info.Subscribe != nil {
				data.Subscribe = info.Subscribe.Name
			}
			// add new
			mapResult[item.SubscribeId] = data
		}
	}
	// convert map to slice
	result = make([]types.ServerOnlineUser, 0, len(mapResult))
	for _, item := range mapResult {
		result = append(result, item)
	}
	return result
}

func (l *FilterServerListLogic) handlerServerStaus(last *time.Time) string {
	if last == nil {
		return "offline"
	}
	if time.Since(*last) > time.Minute*5 {
		return "offline"
	}
	if time.Since(*last) > time.Minute*3 {
		return "warning"
	}
	return "online"

}
