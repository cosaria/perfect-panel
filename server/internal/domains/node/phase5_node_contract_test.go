package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	nodeModel "github.com/perfect-panel/server/internal/platform/persistence/node"
	"github.com/perfect-panel/server/internal/platform/support/tool"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

type failingNodeModel struct {
	nodeModel.Model
}

func (f failingNodeModel) FindOneServer(_ context.Context, _ int64) (*nodeModel.Server, error) {
	return nil, errors.New("lookup failed")
}

func decodeNodeProblem(t *testing.T, recorder *httptest.ResponseRecorder) map[string]interface{} {
	t.Helper()

	var body map[string]interface{}
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &body))
	return body
}

func TestPhase5GetServerConfigReturnsEmpty304OnMatchingETag(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mr := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})

	cached := `{"basic":{"push_interval":60,"pull_interval":60},"protocol":"vless","config":{"port":443}}`
	cacheKey := nodeModel.ServerConfigCacheKey + "1:vless"
	require.NoError(t, client.Set(t.Context(), cacheKey, cached, 0).Err())
	etag := tool.GenerateETag([]byte(cached))

	router := gin.New()
	router.GET("/node/config", GetServerConfigHandler(Deps{Redis: client}))

	req := httptest.NewRequest(http.MethodGet, "/node/config?server_id=1&protocol=vless", nil)
	req.Header.Set("If-None-Match", etag)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusNotModified, resp.Code)
	require.Empty(t, resp.Body.String())
}

func TestPhase5GetServerUserListReturnsCoarseProblemOnLookupFailure(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mr := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})

	router := gin.New()
	router.GET("/node/users", GetServerUserListHandler(Deps{
		Redis:     client,
		NodeModel: failingNodeModel{},
	}))

	req := httptest.NewRequest(http.MethodGet, "/node/users?server_id=1&protocol=vless", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusBadGateway, resp.Code)
	require.Contains(t, resp.Header().Get("Content-Type"), "application/problem+json")

	body := decodeNodeProblem(t, resp)
	require.Equal(t, "Bad Gateway", body["title"])
	require.Equal(t, "Node resource unavailable", body["detail"])
	require.Equal(t, "urn:perfect-panel:error:node-unavailable", body["type"])
	require.NotContains(t, body, "code")
}
