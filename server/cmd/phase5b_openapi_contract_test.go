package cmd_test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPhase5BRepresentativeOperationsExposeGovernedMetadata(t *testing.T) {
	specs := exportPhase5BSpecs(t)

	adminCreateAds := lookupPhase5BOperation(t, specs["admin"], "/ads", "post")
	require.Equal(t, "createAds", adminCreateAds["operationId"])
	require.Equal(t, "Create Ads", adminCreateAds["summary"])
	require.Contains(t, adminCreateAds, "security")
	require.Contains(t, adminCreateAds, "responses")
	require.Contains(t, adminCreateAds["responses"].(map[string]interface{}), "400")
	require.Contains(t, adminCreateAds["responses"].(map[string]interface{}), "401")

	userLogin := lookupPhase5BOperation(t, specs["user"], "/api/v1/auth/login", "post")
	require.Equal(t, "userLogin", userLogin["operationId"])
	require.Equal(t, "User login", userLogin["summary"])
	require.Contains(t, userLogin, "security")
	require.Empty(t, userLogin["security"])
	require.Contains(t, userLogin["responses"].(map[string]interface{}), "400")

	userPurchase := lookupPhase5BOperation(t, specs["user"], "/api/v1/public/order/purchase", "post")
	require.Equal(t, "purchase", userPurchase["operationId"])
	require.Equal(t, "purchase Subscription", userPurchase["summary"])
	require.Contains(t, userPurchase, "security")
	require.Contains(t, userPurchase["responses"].(map[string]interface{}), "401")
}

func lookupPhase5BOperation(t *testing.T, spec map[string]interface{}, path string, method string) map[string]interface{} {
	t.Helper()

	paths, ok := spec["paths"].(map[string]interface{})
	require.True(t, ok)

	pathItem, ok := paths[path].(map[string]interface{})
	require.True(t, ok, "expected path %s in spec", path)

	op, ok := pathItem[method].(map[string]interface{})
	require.True(t, ok, "expected method %s for path %s", method, path)
	return op
}
