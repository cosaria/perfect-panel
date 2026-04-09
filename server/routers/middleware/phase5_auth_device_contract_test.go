package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/perfect-panel/server/config"
	appruntime "github.com/perfect-panel/server/internal/bootstrap/runtime"
	authjwt "github.com/perfect-panel/server/modules/auth/jwt"
	"github.com/stretchr/testify/require"
)

func decodeMiddlewareProblem(t *testing.T, recorder *httptest.ResponseRecorder) map[string]interface{} {
	t.Helper()

	var body map[string]interface{}
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &body))
	return body
}

func TestPhase5AuthMiddlewareCollapsesMissingTokenToCoarseUnauthorizedProblem(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(AuthMiddleware(&appruntime.Deps{
		Config: &config.Config{
			JwtAuth: config.JwtAuth{AccessSecret: "top-secret"},
		},
	}))
	router.GET("/api/v1/user", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/user", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusUnauthorized, resp.Code)
	require.Contains(t, resp.Header().Get("Content-Type"), "application/problem+json")

	body := decodeMiddlewareProblem(t, resp)
	require.Equal(t, "Unauthorized", body["title"])
	require.Equal(t, "Unauthorized", body["detail"])
	require.NotContains(t, body, "code")
}

func TestPhase5AuthMiddlewareCollapsesInvalidTokenToCoarseUnauthorizedProblem(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(AuthMiddleware(&appruntime.Deps{
		Config: &config.Config{
			JwtAuth: config.JwtAuth{AccessSecret: "top-secret"},
		},
	}))
	router.GET("/api/v1/user", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/user", nil)
	req.Header.Set("Authorization", "Bearer definitely-invalid")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusUnauthorized, resp.Code)
	require.Contains(t, resp.Header().Get("Content-Type"), "application/problem+json")

	body := decodeMiddlewareProblem(t, resp)
	require.Equal(t, "Unauthorized", body["title"])
	require.Equal(t, "Unauthorized", body["detail"])
	require.NotContains(t, body, "code")
}

func TestPhase5DeviceMiddlewareCollapsesInvalidCiphertextToBadRequestProblem(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(DeviceMiddleware(&appruntime.Deps{
		Config: &config.Config{
			Device: config.DeviceConfig{
				Enable:         true,
				SecuritySecret: "top-secret",
			},
		},
	}))
	router.POST("/api/v1/device", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/device", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Login-Type", "device")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusBadRequest, resp.Code)
	require.Contains(t, resp.Header().Get("Content-Type"), "application/problem+json")

	body := decodeMiddlewareProblem(t, resp)
	require.Equal(t, "Bad Request", body["title"])
	require.Equal(t, "Invalid request", body["detail"])
	require.NotContains(t, body, "code")
}

func TestPhase5AuthMiddlewareMalformedClaimsDoNotPanic(t *testing.T) {
	gin.SetMode(gin.TestMode)

	token, err := authjwt.NewJwtToken(
		"top-secret",
		time.Now().Unix(),
		3600,
		authjwt.WithOption("SessionId", "session-1"),
	)
	require.NoError(t, err)

	router := gin.New()
	router.Use(AuthMiddleware(&appruntime.Deps{
		Config: &config.Config{
			JwtAuth: config.JwtAuth{AccessSecret: "top-secret"},
		},
	}))
	router.GET("/api/v1/user", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/user", nil)
	req.Header.Set("Authorization", token)
	resp := httptest.NewRecorder()

	require.NotPanics(t, func() {
		router.ServeHTTP(resp, req)
	})

	require.Equal(t, http.StatusUnauthorized, resp.Code)
	require.Contains(t, resp.Header().Get("Content-Type"), "application/problem+json")

	body := decodeMiddlewareProblem(t, resp)
	require.Equal(t, "Unauthorized", body["title"])
	require.Equal(t, "Unauthorized", body["detail"])
	require.NotContains(t, body, "code")
}
