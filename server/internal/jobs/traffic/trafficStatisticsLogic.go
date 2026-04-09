package traffic

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/perfect-panel/server/internal/platform/persistence/node"
	"github.com/perfect-panel/server/internal/platform/support/logger"

	"github.com/hibiken/asynq"
	"github.com/perfect-panel/server/internal/jobs/spec"
	"github.com/perfect-panel/server/internal/platform/persistence/traffic"
)

//goland:noinspection GoNameStartsWithPackageName
type TrafficStatisticsLogic struct {
	deps Deps
}

func NewTrafficStatisticsLogic(deps Deps) *TrafficStatisticsLogic {
	return &TrafficStatisticsLogic{
		deps: deps,
	}
}

func (l *TrafficStatisticsLogic) ProcessTask(ctx context.Context, task *asynq.Task) error {
	var payload spec.TrafficStatistics
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		logger.WithContext(ctx).Error("[TrafficStatistics] Unmarshal payload failed",
			logger.Field("error", err.Error()),
			logger.Field("payload", string(task.Payload())),
		)
		return nil
	}
	if len(payload.Logs) == 0 {
		logger.WithContext(ctx).Error("[TrafficStatistics] Payload is empty")
		return nil
	}
	// query server info
	serverInfo, err := l.deps.NodeModel.FindOneServer(ctx, payload.ServerId)
	if err != nil {
		logger.WithContext(ctx).Error("[TrafficStatistics] Find server info failed",
			logger.Field("serverId", payload.ServerId),
			logger.Field("error", err.Error()),
		)
		return nil
	}
	// query protocol ratio
	// default ratio is 1.0

	protocols, err := serverInfo.UnmarshalProtocols()
	if err != nil {
		logger.Errorf("[TrafficStatistics] Unmarshal protocols failed: %s", err.Error())
		return nil
	}
	var protocol *node.Protocol

	var ratio float32 = 1.0

	for _, p := range protocols {
		if strings.EqualFold(p.Type, payload.Protocol) {
			protocol = &p
			break
		}
	}

	if protocol == nil {
		logger.WithContext(ctx).Error("[TrafficStatistics] Protocol not found: %s", payload.Protocol)
		return nil
	}

	// use protocol ratio if it's greater than 0
	if protocol.Ratio > 0 {
		ratio = float32(protocol.Ratio)
	}

	now := time.Now()
	realTimeMultiplier := float32(1.0)
	manager, err := l.deps.ResolveNodeMultiplierManager(ctx)
	if err != nil {
		logger.WithContext(ctx).Error("[TrafficStatisticsLogic] Resolve node multiplier manager failed",
			logger.Field("error", err.Error()),
		)
		return err
	}
	if manager != nil {
		realTimeMultiplier = manager.GetMultiplier(now)
	}
	cfg := l.deps.currentConfig()
	logger.Debugf("[TrafficStatisticsLogic] Current time traffic multiplier: %.2f", realTimeMultiplier)
	for _, log := range payload.Logs {
		// query user Subscribe Info
		sub, err := l.deps.UserModel.FindOneSubscribe(ctx, log.SID)
		if err != nil {
			logger.WithContext(ctx).Error("[TrafficStatistics] Find user Subscribe Info failed",
				logger.Field("uid", log.SID),
				logger.Field("error", err.Error()),
			)
			continue
		}

		if log.Download+log.Upload <= cfg.Node.TrafficReportThreshold {
			// no traffic, skip
			continue
		}
		// update user subscribe with log
		d := int64(float32(log.Download) * ratio * realTimeMultiplier)
		u := int64(float32(log.Upload) * ratio * realTimeMultiplier)
		if err := l.deps.UserModel.UpdateUserSubscribeWithTraffic(ctx, sub.Id, d, u); err != nil {
			logger.WithContext(ctx).Error("[TrafficStatistics] Update user subscribe with log failed",
				logger.Field("sid", log.SID),
				logger.Field("download", float32(log.Download)*ratio),
				logger.Field("upload", float32(log.Upload)*ratio),
				logger.Field("error", err.Error()),
			)
			continue
		}

		// create log log
		if err = l.deps.TrafficLogModel.Insert(ctx, &traffic.TrafficLog{
			ServerId:    payload.ServerId,
			SubscribeId: log.SID,
			UserId:      sub.UserId,
			Upload:      u,
			Download:    d,
			Timestamp:   now,
		}); err != nil {
			logger.WithContext(ctx).Error("[TrafficStatistics] Create log log failed",
				logger.Field("uid", log.SID),
				logger.Field("download", float32(log.Download)*ratio),
				logger.Field("upload", float32(log.Upload)*ratio),
				logger.Field("error", err.Error()),
			)
		}
	}
	return nil
}
