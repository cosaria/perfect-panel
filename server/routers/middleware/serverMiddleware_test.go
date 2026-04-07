package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/svc"
	"github.com/stretchr/testify/require"
)

func TestServerMiddlewareAcceptsHeaderSecret(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(ServerMiddleware(&svc.ServiceContext{
		Config: config.Config{
			Node: config.NodeConfig{
				NodeSecret: "top-secret",
			},
		},
	}))
	router.GET("/node", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/node", nil)
	req.Header.Set("X-Node-Secret", "top-secret")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
	require.Equal(t, "ok", resp.Body.String())
}

func TestServerMiddlewareRejectsQuerySecret(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(ServerMiddleware(&svc.ServiceContext{
		Config: config.Config{
			Node: config.NodeConfig{
				NodeSecret: "top-secret",
			},
		},
	}))
	router.GET("/node", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/node?secret_key=top-secret", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusForbidden, resp.Code)
	require.Equal(t, "Forbidden", resp.Body.String())
}
