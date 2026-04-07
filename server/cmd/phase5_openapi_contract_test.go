package cmd_test

import (
	"encoding/json"
	"testing"

	"github.com/gin-gonic/gin"
	handler "github.com/perfect-panel/server/routers"
	"github.com/stretchr/testify/require"
)

func TestPhase5OpenAPIExportsProblemContract(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	apis := handler.RegisterHandlersForSpec(router)

	userSpec, err := apis.UserOpenAPI()
	require.NoError(t, err)

	assertSpecHasProblemContract(t, apis.Admin.OpenAPI())
	assertSpecHasProblemContract(t, apis.Common.OpenAPI())
	assertSpecHasProblemContract(t, userSpec)
}

func assertSpecHasProblemContract(t *testing.T, spec interface{}) {
	t.Helper()

	data, err := json.Marshal(spec)
	require.NoError(t, err)

	var body map[string]interface{}
	require.NoError(t, json.Unmarshal(data, &body))

	components, ok := body["components"].(map[string]interface{})
	require.True(t, ok)

	schemas, ok := components["schemas"].(map[string]interface{})
	require.True(t, ok)
	require.Contains(t, schemas, "Problem")

	specJSON := string(data)
	require.Contains(t, specJSON, "\"application/problem+json\"")
	require.Contains(t, specJSON, "\"#/components/schemas/Problem\"")
}
