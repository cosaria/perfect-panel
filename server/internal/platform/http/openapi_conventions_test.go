package handler

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGovernedAPIConfigIncludesSecuritySchemeAndTags(t *testing.T) {
	cfg := governedAPIConfig("Test API", "1.0.0", "/api/v1/test", "auth", "oauth")

	require.Len(t, cfg.Servers, 1)
	require.Equal(t, "/api/v1/test", cfg.Servers[0].URL)
	require.Contains(t, cfg.Components.SecuritySchemes, "bearer")
	require.Len(t, cfg.Tags, 2)
	require.Equal(t, "auth", cfg.Tags[0].Name)
	require.NotEmpty(t, cfg.Tags[0].Description)
	require.Equal(t, "oauth", cfg.Tags[1].Name)
}
