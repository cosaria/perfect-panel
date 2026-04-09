package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	appruntime "github.com/perfect-panel/server/internal/bootstrap/runtime"
	"github.com/perfect-panel/server/internal/domains/subscribe"
	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/util/tool"
)

func PanDomainMiddleware(runtimeDeps *appruntime.Deps) func(c *gin.Context) {
	return func(c *gin.Context) {

		if runtimeDeps != nil && runtimeDeps.Config != nil && runtimeDeps.Config.Subscribe.PanDomain && c.Request.URL.Path == "/" {
			// intercept browser
			ua := c.GetHeader("User-Agent")

			if runtimeDeps.Config.Subscribe.UserAgentLimit {
				if ua == "" {
					c.String(http.StatusForbidden, "Access denied")
					c.Abort()
					return
				}
				browserKeywords := tool.RemoveDuplicateElements(strings.Split(runtimeDeps.Config.Subscribe.UserAgentList, "\n")...)
				var allow = false

				// query client list
				clients, err := runtimeDeps.ClientModel.List(c.Request.Context())
				if err != nil {
					logger.Errorw("[PanDomainMiddleware] Query client list failed", logger.Field("error", err.Error()))
				}
				for _, item := range clients {
					u := strings.ToLower(item.UserAgent)
					u = strings.Trim(u, " ")
					browserKeywords = append(browserKeywords, u)
				}

				for _, keyword := range browserKeywords {
					keyword = strings.ToLower(strings.Trim(keyword, " "))
					if keyword == "" {
						continue
					}
					if strings.Contains(strings.ToLower(ua), keyword) {
						allow = true
					}
				}
				if !allow {
					c.String(http.StatusForbidden, "Access denied")
					c.Abort()
					return
				}
			}

			domain := c.Request.Host
			domainArr := strings.Split(domain, ".")
			domainFirst := domainArr[0]
			request := types.SubscribeRequest{
				Token: domainFirst,
				Flag:  domainArr[1],
				UA:    c.Request.Header.Get("User-Agent"),
			}
			l := subscribe.NewSubscribeLogic(c, subscribe.Deps{
				ClientModel:    runtimeDeps.ClientModel,
				LogModel:       runtimeDeps.LogModel,
				NodeModel:      runtimeDeps.NodeModel,
				SubscribeModel: runtimeDeps.SubscribeModel,
				UserModel:      runtimeDeps.UserModel,
				Config:         runtimeDeps.Config,
			})
			resp, err := l.Handler(&request)
			if err != nil {
				return
			}
			c.Header("subscription-userinfo", resp.Header)
			c.String(200, "%s", string(resp.Config))
			c.Abort()
			return
		}
		c.Next()
	}
}
