package cmd_test

import (
	"encoding/json"
	"testing"

	"github.com/gin-gonic/gin"
	handler "github.com/perfect-panel/server/internal/platform/http"
	"github.com/stretchr/testify/require"
)

func exportPhase5BSpecs(t *testing.T) map[string]map[string]interface{} {
	t.Helper()

	gin.SetMode(gin.TestMode)

	router := gin.New()
	apis := handler.RegisterHandlersForSpec(router)

	userSpec, err := apis.UserOpenAPI()
	require.NoError(t, err)

	return map[string]map[string]interface{}{
		"admin":  decodePhase5BSpec(t, apis.Admin.OpenAPI()),
		"common": decodePhase5BSpec(t, apis.Common.OpenAPI()),
		"user":   decodePhase5BSpec(t, userSpec),
	}
}

func decodePhase5BSpec(t *testing.T, spec interface{}) map[string]interface{} {
	t.Helper()

	data, err := json.Marshal(spec)
	require.NoError(t, err)

	var body map[string]interface{}
	require.NoError(t, json.Unmarshal(data, &body))
	return body
}
