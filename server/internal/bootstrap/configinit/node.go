package configinit

import (
	"context"
	"encoding/json"

	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/models/node"
	"github.com/perfect-panel/server/models/system"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/util/tool"
)

func Node(deps Deps) {
	logger.Debug("Node config initialization")
	configs, err := deps.SystemModel.GetNodeConfig(context.Background())
	if err != nil {
		panic(err)
	}
	var nodeConfig config.NodeDBConfig
	tool.SystemConfigSliceReflectToStruct(configs, &nodeConfig)
	c := config.NodeConfig{
		NodeSecret:             nodeConfig.NodeSecret,
		NodePullInterval:       nodeConfig.NodePullInterval,
		NodePushInterval:       nodeConfig.NodePushInterval,
		IPStrategy:             nodeConfig.IPStrategy,
		TrafficReportThreshold: nodeConfig.TrafficReportThreshold,
	}
	if nodeConfig.DNS != "" {
		var dns []config.NodeDNS
		err = json.Unmarshal([]byte(nodeConfig.DNS), &dns)
		if err != nil {
			logger.Errorf("[Node] Unmarshal DNS config error: %s", err.Error())
			panic(err)
		}
		c.DNS = dns
	}
	if nodeConfig.Block != "" {
		var block []string
		_ = json.Unmarshal([]byte(nodeConfig.Block), &block)
		c.Block = tool.RemoveDuplicateElements(block...)
	}
	if nodeConfig.Outbound != "" {
		var outbound []config.NodeOutbound
		err = json.Unmarshal([]byte(nodeConfig.Outbound), &outbound)
		if err != nil {
			logger.Errorf("[Node] Unmarshal Outbound config error: %s", err.Error())
			panic(err)
		}
		c.Outbound = outbound
	}

	if deps.Config != nil {
		deps.Config.Node = c
	}

	// Manager initialization
	if deps.DB.Model(&system.System{}).Where("`key` = ?", "NodeMultiplierConfig").Find(&system.System{}).RowsAffected == 0 {
		if err := deps.DB.Model(&system.System{}).Create(&system.System{
			Key:      "NodeMultiplierConfig",
			Value:    "[]",
			Type:     "string",
			Desc:     "Node Multiplier Config",
			Category: "server",
		}).Error; err != nil {
			logger.Errorf("Create Node Multiplier Config Error: %s", err.Error())
		}
		return
	}

	nodeMultiplierData, err := deps.SystemModel.FindNodeMultiplierConfig(context.Background())
	if err != nil {
		logger.Error("Get Node Multiplier Config Error: ", logger.Field("error", err.Error()))
		return
	}
	var periods []node.TimePeriod
	if err := json.Unmarshal([]byte(nodeMultiplierData.Value), &periods); err != nil {
		logger.Error("Unmarshal Node Multiplier Config Error: ", logger.Field("error", err.Error()), logger.Field("value", nodeMultiplierData.Value))
	}
	if deps.SetNodeMultiplierManager != nil {
		deps.SetNodeMultiplierManager(node.NewNodeMultiplierManager(periods))
	}
}
