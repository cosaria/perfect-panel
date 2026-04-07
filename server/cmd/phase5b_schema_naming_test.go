package cmd_test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPhase5BSchemaNamesAvoidKnownTypos(t *testing.T) {
	specs := exportPhase5BSpecs(t)

	userSchemas := phase5BSchemas(t, specs["user"])
	require.Contains(t, userSchemas, "OAuthLoginRequest")
	require.NotContains(t, userSchemas, "OAthLoginRequest")

	adminSchemas := phase5BSchemas(t, specs["admin"])
	require.Contains(t, adminSchemas, "DeleteUserDeviceRequest")
	require.NotContains(t, adminSchemas, "DeleteUserDeivceRequest")
}

func phase5BSchemas(t *testing.T, spec map[string]interface{}) map[string]interface{} {
	t.Helper()

	components, ok := spec["components"].(map[string]interface{})
	require.True(t, ok)

	schemas, ok := components["schemas"].(map[string]interface{})
	require.True(t, ok)
	return schemas
}
