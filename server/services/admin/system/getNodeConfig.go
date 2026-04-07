package system

import (
	"context"
	"encoding/json"

	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type GetNodeConfigOutput struct {
	Body *types.NodeConfig
}

func GetNodeConfigHandler(deps Deps) func(context.Context, *struct{}) (*GetNodeConfigOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetNodeConfigOutput, error) {
		l := NewGetNodeConfigLogic(ctx, deps)
		resp, err := l.GetNodeConfig()
		if err != nil {
			return nil, err
		}
		return &GetNodeConfigOutput{Body: resp}, nil
	}
}

type GetNodeConfigLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

func NewGetNodeConfigLogic(ctx context.Context, deps Deps) *GetNodeConfigLogic {
	return &GetNodeConfigLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *GetNodeConfigLogic) GetNodeConfig() (*types.NodeConfig, error) {
	// get server config from db
	configs, err := l.deps.SystemModel.GetNodeConfig(l.ctx)
	if err != nil {
		l.Errorw("[GetNodeConfigLogic] GetNodeConfig get server config error: ", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "GetNodeConfig get server config error: %v", err.Error())
	}
	var dbConfig config.NodeDBConfig
	tool.SystemConfigSliceReflectToStruct(configs, &dbConfig)
	c := &types.NodeConfig{
		NodeSecret:             dbConfig.NodeSecret,
		NodePullInterval:       dbConfig.NodePullInterval,
		NodePushInterval:       dbConfig.NodePushInterval,
		IPStrategy:             dbConfig.IPStrategy,
		TrafficReportThreshold: dbConfig.TrafficReportThreshold,
	}

	if dbConfig.DNS != "" {
		var dns []types.NodeDNS
		err = json.Unmarshal([]byte(dbConfig.DNS), &dns)
		if err != nil {
			logger.Errorf("[Node] Unmarshal DNS config error: %s", err.Error())
			panic(err)
		}
		c.DNS = dns
	}
	if dbConfig.Block != "" {
		var block []string
		_ = json.Unmarshal([]byte(dbConfig.Block), &block)
		c.Block = tool.RemoveDuplicateElements(block...)
	}
	if dbConfig.Outbound != "" {
		var outbound []types.NodeOutbound
		err = json.Unmarshal([]byte(dbConfig.Outbound), &outbound)
		if err != nil {
			logger.Errorf("[Node] Unmarshal Outbound config error: %s", err.Error())
			panic(err)
		}
		c.Outbound = outbound
	}

	return c, nil
}
