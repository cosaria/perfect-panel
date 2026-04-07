package subscribe

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/perfect-panel/server/adapter"
	"github.com/perfect-panel/server/models/client"
	"github.com/perfect-panel/server/models/log"
	"github.com/perfect-panel/server/models/node"
	"github.com/perfect-panel/server/models/user"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/services/report"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func SubscribeHandler(svcCtx *svc.ServiceContext) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req types.SubscribeRequest
		if c.Request.Header.Get("token") != "" {
			req.Token = c.Request.Header.Get("token")
		} else {
			req.Token = c.Query("token")
		}
		ua := c.GetHeader("User-Agent")
		req.UA = c.Request.Header.Get("User-Agent")
		req.Flag = c.Query("flag")
		req.Type = c.Query("type")
		// 获取所有查询参数
		req.Params = getQueryMap(c.Request)

		if svcCtx.Config.Subscribe.PanDomain {
			domain := c.Request.Host
			domainArr := strings.Split(domain, ".")
			short, err := tool.FixedUniqueString(req.Token, 8, "")
			if err != nil {
				logger.Errorf("[SubscribeHandler] Generate short token failed: %v", err)
				c.String(http.StatusInternalServerError, "Internal Server")
				c.Abort()
				return
			}
			if !strings.EqualFold(short, domainArr[0]) {
				logger.Debugf("[SubscribeHandler] Generate short token failed, short: %s, domain: %s", short, domainArr[0])
				c.String(http.StatusForbidden, "Access denied")
				c.Abort()
				return
			}
		}

		if svcCtx.Config.Subscribe.UserAgentLimit {
			if ua == "" {
				c.String(http.StatusForbidden, "Access denied")
				c.Abort()
				return
			}
			clientUserAgents := tool.RemoveDuplicateElements(strings.Split(svcCtx.Config.Subscribe.UserAgentList, "\n")...)

			// query client list
			clients, err := svcCtx.ClientModel.List(c.Request.Context())
			if err != nil {
				logger.Errorw("[PanDomainMiddleware] Query client list failed", logger.Field("error", err.Error()))
			}
			for _, item := range clients {
				u := strings.ToLower(item.UserAgent)
				u = strings.Trim(u, " ")
				clientUserAgents = append(clientUserAgents, u)
			}

			var allow = false
			for _, keyword := range clientUserAgents {
				keyword = strings.Trim(keyword, " ")
				if keyword == "" {
					continue
				}
				if strings.Contains(strings.ToLower(ua), strings.ToLower(keyword)) {
					allow = true
				}
			}
			if !allow {
				c.String(http.StatusForbidden, "Access denied")
				c.Abort()
				return
			}
		}

		l := NewSubscribeLogic(c, svcCtx)
		resp, err := l.Handler(&req)
		if err != nil {
			c.String(http.StatusInternalServerError, "Internal Server")
			return
		}
		c.Header("subscription-userinfo", resp.Header)
		c.String(200, "%s", string(resp.Config))
	}
}

func RegisterSubscribeHandlers(router *gin.Engine, serverCtx *svc.ServiceContext) {
	path := serverCtx.Config.Subscribe.SubscribePath
	if path == "" {
		path = "/api/v1/subscribe/config"
	}
	router.GET(path, SubscribeHandler(serverCtx))
}

// GetQueryMap 将 http.Request 的查询参数转换为 map[string]string
func getQueryMap(r *http.Request) map[string]string {
	result := make(map[string]string)
	for k, v := range r.URL.Query() {
		if len(v) > 0 {
			result[k] = v[0]
		}
	}
	return result
}

type SubscribeLogic struct {
	ctx *gin.Context
	svc *svc.ServiceContext
	logger.Logger
}

func NewSubscribeLogic(ctx *gin.Context, svc *svc.ServiceContext) *SubscribeLogic {
	return &SubscribeLogic{
		ctx:    ctx,
		svc:    svc,
		Logger: logger.WithContext(ctx.Request.Context()),
	}
}

func (l *SubscribeLogic) Handler(req *types.SubscribeRequest) (resp *types.SubscribeResponse, err error) {
	// query client list
	clients, err := l.svc.ClientModel.List(l.ctx.Request.Context())
	if err != nil {
		l.Errorw("[SubscribeLogic] Query client list failed", logger.Field("error", err.Error()))
		return nil, err
	}

	userAgent := strings.ToLower(l.ctx.Request.UserAgent())

	var targetApp, defaultApp *client.SubscribeApplication

	for _, item := range clients {
		u := strings.ToLower(item.UserAgent)
		if item.IsDefault {
			defaultApp = item
		}

		if strings.Contains(userAgent, u) {
			// Special handling for Stash
			if strings.Contains(userAgent, "stash") && !strings.Contains(u, "stash") {
				continue
			}
			targetApp = item
			break
		}
	}
	if targetApp == nil {
		l.Debugf("[SubscribeLogic] No matching client found", logger.Field("userAgent", userAgent))
		if defaultApp == nil {
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "No matching client found for user agent: %s", userAgent)
		}
		targetApp = defaultApp
	}
	// Find user subscribe by token
	userSubscribe, err := l.getUserSubscribe(req.Token)
	if err != nil {
		l.Errorw("[SubscribeLogic] Get user subscribe failed", logger.Field("error", err.Error()), logger.Field("token", req.Token))
		return nil, err
	}

	var subscribeStatus = false
	defer func() {
		l.logSubscribeActivity(subscribeStatus, userSubscribe, req)
	}()
	// find subscribe info
	subscribeInfo, err := l.svc.SubscribeModel.FindOne(l.ctx.Request.Context(), userSubscribe.SubscribeId)
	if err != nil {
		l.Errorw("[SubscribeLogic] Find subscribe info failed", logger.Field("error", err.Error()), logger.Field("subscribeId", userSubscribe.SubscribeId))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "Find subscribe info failed: %v", err.Error())
	}

	// Find server list by user subscribe
	servers, err := l.getServers(userSubscribe)
	if err != nil {
		return nil, err
	}
	a := adapter.NewAdapter(
		targetApp.SubscribeTemplate,
		adapter.WithServers(servers),
		adapter.WithSiteName(l.svc.Config.Site.SiteName),
		adapter.WithSubscribeName(subscribeInfo.Name),
		adapter.WithOutputFormat(targetApp.OutputFormat),
		adapter.WithUserInfo(adapter.User{
			Password:     userSubscribe.UUID,
			ExpiredAt:    userSubscribe.ExpireTime,
			Download:     userSubscribe.Download,
			Upload:       userSubscribe.Upload,
			Traffic:      userSubscribe.Traffic,
			SubscribeURL: l.getSubscribeV2URL(),
		}),
		adapter.WithParams(req.Params),
	)

	logger.Debugf("[SubscribeLogic] Building client config for user %d with URI %s", userSubscribe.UserId, l.getSubscribeV2URL())

	// Get client config
	adapterClient, err := a.Client()
	if err != nil {
		l.Errorw("[SubscribeLogic] Client error", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(500), "Client error: %v", err.Error())
	}
	bytes, err := adapterClient.Build()
	if err != nil {
		l.Errorw("[SubscribeLogic] Build client config failed", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(500), "Build client config failed: %v", err.Error())
	}

	var formats = []string{"json", "yaml", "conf"}

	for _, format := range formats {
		if format == strings.ToLower(targetApp.OutputFormat) {
			l.ctx.Header("content-disposition", fmt.Sprintf("attachment;filename*=UTF-8''%s.%s", url.QueryEscape(l.svc.Config.Site.SiteName), format))
			l.ctx.Header("Content-Type", "application/octet-stream; charset=UTF-8")

		}
	}

	resp = &types.SubscribeResponse{
		Config: bytes,
		Header: fmt.Sprintf(
			"upload=%d;download=%d;total=%d;expire=%d",
			userSubscribe.Upload, userSubscribe.Download, userSubscribe.Traffic, userSubscribe.ExpireTime.Unix(),
		),
	}
	subscribeStatus = true
	return
}

func (l *SubscribeLogic) getSubscribeV2URL() string {

	uri := l.ctx.Request.RequestURI
	// is gateway mode, add /sub prefix
	if report.IsGatewayMode() {
		uri = "/sub" + uri
	}
	// use custom domain if configured
	if l.svc.Config.Subscribe.SubscribeDomain != "" {
		domains := strings.Split(l.svc.Config.Subscribe.SubscribeDomain, "\n")
		return fmt.Sprintf("https://%s%s", domains[0], uri)
	}
	// use current request host
	return fmt.Sprintf("https://%s%s", l.ctx.Request.Host, uri)
}

func (l *SubscribeLogic) getUserSubscribe(token string) (*user.Subscribe, error) {
	userSub, err := l.svc.UserModel.FindOneSubscribeByToken(l.ctx.Request.Context(), token)
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

func (l *SubscribeLogic) logSubscribeActivity(subscribeStatus bool, userSub *user.Subscribe, req *types.SubscribeRequest) {
	if !subscribeStatus {
		return
	}

	subscribeLog := log.Subscribe{
		Token:           req.Token,
		UserAgent:       req.UA,
		ClientIP:        l.ctx.ClientIP(),
		UserSubscribeId: userSub.Id,
	}

	content, _ := subscribeLog.Marshal()

	err := l.svc.LogModel.Insert(l.ctx.Request.Context(), &log.SystemLog{
		Type:     log.TypeSubscribe.Uint8(),
		ObjectID: userSub.UserId, // log user id
		Date:     time.Now().Format(time.DateOnly),
		Content:  string(content),
	})
	if err != nil {
		l.Errorw("[Generate Subscribe]insert subscribe log error: %v", logger.Field("error", err.Error()))
	}
}

func (l *SubscribeLogic) getServers(userSub *user.Subscribe) ([]*node.Node, error) {
	if l.isSubscriptionExpired(userSub) {
		return l.createExpiredServers(), nil
	}

	subDetails, err := l.svc.SubscribeModel.FindOne(l.ctx.Request.Context(), userSub.SubscribeId)
	if err != nil {
		l.Errorw("[Generate Subscribe]find subscribe details error: %v", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "find subscribe details error: %v", err.Error())
	}

	nodeIds := tool.StringToInt64Slice(subDetails.Nodes)
	tags := tool.RemoveStringElement(strings.Split(subDetails.NodeTags, ","), "")

	l.Debugf("[Generate Subscribe]nodes: %v, NodeTags: %v", len(nodeIds), len(tags))
	if len(nodeIds) == 0 && len(tags) == 0 {
		logger.Infow("[Generate Subscribe]no subscribe nodes")
		return []*node.Node{}, nil
	}
	enable := true
	var nodes []*node.Node
	_, nodes, err = l.svc.NodeModel.FilterNodeList(l.ctx.Request.Context(), &node.FilterNodeParams{
		Page:    1,
		Size:    1000,
		NodeId:  nodeIds,
		Tag:     tool.RemoveDuplicateElements(tags...),
		Preload: true,
		Enabled: &enable, // Only get enabled nodes
	})

	l.Debugf("[Query Subscribe]found servers: %v", len(nodes))

	if err != nil {
		l.Errorw("[Generate Subscribe]find server details error: %v", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "find server details error: %v", err.Error())
	}
	logger.Debugf("[Generate Subscribe]found servers: %v", len(nodes))
	return nodes, nil
}

func (l *SubscribeLogic) isSubscriptionExpired(userSub *user.Subscribe) bool {
	return userSub.ExpireTime.Unix() < time.Now().Unix() && userSub.ExpireTime.Unix() != 0
}

func (l *SubscribeLogic) createExpiredServers() []*node.Node {
	enable := true
	host := l.getFirstHostLine()

	return []*node.Node{
		{
			Name:    "Subscribe Expired",
			Tags:    "",
			Port:    18080,
			Address: "127.0.0.1",
			Server: &node.Server{
				Id:        1,
				Name:      "Subscribe Expired",
				Protocols: "[{\"type\":\"shadowsocks\",\"cipher\":\"aes-256-gcm\",\"port\":1}]",
			},
			Protocol: "shadowsocks",
			Enabled:  &enable,
		},
		{
			Name:    host,
			Tags:    "",
			Port:    18080,
			Address: "127.0.0.1",
			Server: &node.Server{
				Id:        1,
				Name:      "Subscribe Expired",
				Protocols: "[{\"type\":\"shadowsocks\",\"cipher\":\"aes-256-gcm\",\"port\":1}]",
			},
			Protocol: "shadowsocks",
			Enabled:  &enable,
		},
	}
}

func (l *SubscribeLogic) getFirstHostLine() string {
	host := l.svc.Config.Host
	lines := strings.Split(host, "\n")
	if len(lines) > 0 {
		return lines[0]
	}
	return host
}
