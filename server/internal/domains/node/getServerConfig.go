package server

import (
	"encoding/base64"
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
)

func GetServerConfigHandler(deps Deps) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req types.GetServerConfigRequest
		_ = c.ShouldBind(&req)
		_ = c.ShouldBindQuery(&req.ServerCommon)
		validateErr := validateRequest(&req)
		if validateErr != nil {
			response.ParamErrorResult(c, validateErr)
			return
		}

		l := NewGetServerConfigLogic(c, deps)
		resp, err := l.GetServerConfig(&req)
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

type GetServerConfigLogic struct {
	logger.Logger
	ctx  *gin.Context
	deps Deps
}

// NewGetServerConfigLogic Get server config
func NewGetServerConfigLogic(ctx *gin.Context, deps Deps) *GetServerConfigLogic {
	return &GetServerConfigLogic{
		Logger: logger.WithContext(ctx.Request.Context()),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *GetServerConfigLogic) GetServerConfig(req *types.GetServerConfigRequest) (resp *types.GetServerConfigResponse, err error) {
	cacheKey := fmt.Sprintf("%s%d:%s", node.ServerConfigCacheKey, req.ServerId, req.Protocol)
	cache, err := l.deps.Redis.Get(l.ctx, cacheKey).Result()
	if err == nil {
		if cache != "" {
			etag := tool.GenerateETag([]byte(cache))
			//  Check If-None-Match header
			match := l.ctx.GetHeader("If-None-Match")
			if match == etag {
				return nil, xerr.ErrStatusNotModified
			}
			l.ctx.Header("ETag", etag)
			resp = &types.GetServerConfigResponse{}
			err = json.Unmarshal([]byte(cache), resp)
			if err != nil {
				l.Errorw("[ServerConfigCacheKey] json unmarshal error", logger.Field("error", err.Error()))
				return nil, err
			}
			return resp, nil
		}
	}
	data, err := l.deps.NodeModel.FindOneServer(l.ctx, req.ServerId)
	if err != nil {
		l.Errorw("[GetServerConfig] FindOne error", logger.Field("error", err.Error()))
		return nil, err
	}

	// compatible hysteria2, remove in future versions
	protocolRequest := req.Protocol
	if protocolRequest == Hysteria2 {
		protocolRequest = Hysteria
	}

	protocols, err := data.UnmarshalProtocols()
	if err != nil {
		return nil, err
	}
	var cfg map[string]interface{}
	for _, protocol := range protocols {
		if protocol.Type == protocolRequest {
			cfg = l.compatible(protocol)
			break
		}
	}

	appCfg := l.deps.currentConfig()
	resp = &types.GetServerConfigResponse{
		Basic: types.ServerBasic{
			PullInterval: appCfg.Node.NodePullInterval,
			PushInterval: appCfg.Node.NodePushInterval,
		},
		Protocol: req.Protocol,
		Config:   cfg,
	}
	c, err := json.Marshal(resp)
	if err != nil {
		l.Errorw("[GetServerConfig] json marshal error", logger.Field("error", err.Error()))
		return nil, err
	}
	etag := tool.GenerateETag(c)
	l.ctx.Header("ETag", etag)
	if err = l.deps.Redis.Set(l.ctx, cacheKey, c, -1).Err(); err != nil {
		l.Errorw("[GetServerConfig] redis set error", logger.Field("error", err.Error()))
	}
	//  Check If-None-Match header
	match := l.ctx.GetHeader("If-None-Match")
	if match == etag {
		return nil, xerr.ErrStatusNotModified
	}

	return resp, nil
}

func (l *GetServerConfigLogic) compatible(config node.Protocol) map[string]interface{} {
	var result interface{}
	switch config.Type {
	case ShadowSocks:
		result = ShadowsocksNode{
			Port:      config.Port,
			Cipher:    config.Cipher,
			ServerKey: base64.StdEncoding.EncodeToString([]byte(config.ServerKey)),
		}
	case Vless:
		result = VlessNode{
			Port:    config.Port,
			Flow:    config.Flow,
			Network: config.Transport,
			TransportConfig: &TransportConfig{
				Path:                 config.Path,
				Host:                 config.Host,
				ServiceName:          config.ServiceName,
				DisableSNI:           config.DisableSNI,
				ReduceRtt:            config.ReduceRtt,
				UDPRelayMode:         config.UDPRelayMode,
				CongestionController: config.CongestionController,
			},
			Security: config.Security,
			SecurityConfig: &SecurityConfig{
				SNI:                  config.SNI,
				AllowInsecure:        &config.AllowInsecure,
				Fingerprint:          config.Fingerprint,
				RealityServerAddress: config.RealityServerAddr,
				RealityServerPort:    config.RealityServerPort,
				RealityPrivateKey:    config.RealityPrivateKey,
				RealityPublicKey:     config.RealityPublicKey,
				RealityShortId:       config.RealityShortId,
			},
		}
	case Vmess:
		result = VmessNode{
			Port:    config.Port,
			Network: config.Transport,
			TransportConfig: &TransportConfig{
				Path:                 config.Path,
				Host:                 config.Host,
				ServiceName:          config.ServiceName,
				DisableSNI:           config.DisableSNI,
				ReduceRtt:            config.ReduceRtt,
				UDPRelayMode:         config.UDPRelayMode,
				CongestionController: config.CongestionController,
			},
			Security: config.Security,
			SecurityConfig: &SecurityConfig{
				SNI:                  config.SNI,
				AllowInsecure:        &config.AllowInsecure,
				Fingerprint:          config.Fingerprint,
				RealityServerAddress: config.RealityServerAddr,
				RealityServerPort:    config.RealityServerPort,
				RealityPrivateKey:    config.RealityPrivateKey,
				RealityPublicKey:     config.RealityPublicKey,
				RealityShortId:       config.RealityShortId,
			},
		}
	case Trojan:
		result = TrojanNode{
			Port:    config.Port,
			Network: config.Transport,
			TransportConfig: &TransportConfig{
				Path:                 config.Path,
				Host:                 config.Host,
				ServiceName:          config.ServiceName,
				DisableSNI:           config.DisableSNI,
				ReduceRtt:            config.ReduceRtt,
				UDPRelayMode:         config.UDPRelayMode,
				CongestionController: config.CongestionController,
			},
			Security: config.Security,
			SecurityConfig: &SecurityConfig{
				SNI:                  config.SNI,
				AllowInsecure:        &config.AllowInsecure,
				Fingerprint:          config.Fingerprint,
				RealityServerAddress: config.RealityServerAddr,
				RealityServerPort:    config.RealityServerPort,
				RealityPrivateKey:    config.RealityPrivateKey,
				RealityPublicKey:     config.RealityPublicKey,
				RealityShortId:       config.RealityShortId,
			},
		}
	case AnyTLS:
		result = AnyTLSNode{
			Port: config.Port,
			SecurityConfig: &SecurityConfig{
				SNI:                  config.SNI,
				AllowInsecure:        &config.AllowInsecure,
				Fingerprint:          config.Fingerprint,
				RealityServerAddress: config.RealityServerAddr,
				RealityServerPort:    config.RealityServerPort,
				RealityPrivateKey:    config.RealityPrivateKey,
				RealityPublicKey:     config.RealityPublicKey,
				RealityShortId:       config.RealityShortId,
				PaddingScheme:        config.PaddingScheme,
			},
		}
	case Tuic:
		result = TuicNode{
			Port: config.Port,
			SecurityConfig: &SecurityConfig{
				SNI:                  config.SNI,
				AllowInsecure:        &config.AllowInsecure,
				Fingerprint:          config.Fingerprint,
				RealityServerAddress: config.RealityServerAddr,
				RealityServerPort:    config.RealityServerPort,
				RealityPrivateKey:    config.RealityPrivateKey,
				RealityPublicKey:     config.RealityPublicKey,
				RealityShortId:       config.RealityShortId,
			},
		}
	case Hysteria:
		result = Hysteria2Node{
			Port:         config.Port,
			HopPorts:     config.HopPorts,
			HopInterval:  config.HopInterval,
			ObfsPassword: config.ObfsPassword,
			SecurityConfig: &SecurityConfig{
				SNI:                  config.SNI,
				AllowInsecure:        &config.AllowInsecure,
				Fingerprint:          config.Fingerprint,
				RealityServerAddress: config.RealityServerAddr,
				RealityServerPort:    config.RealityServerPort,
				RealityPrivateKey:    config.RealityPrivateKey,
				RealityPublicKey:     config.RealityPublicKey,
				RealityShortId:       config.RealityShortId,
			},
		}

	}
	var resp map[string]interface{}
	s, _ := json.Marshal(result)
	_ = json.Unmarshal(s, &resp)
	return resp
}
