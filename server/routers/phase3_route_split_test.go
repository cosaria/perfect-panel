package handler

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestPhase3RouteRegistrarsComposeSpecContract(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	apis := &APIs{}

	registerAdminRoutes(router, nil, true, apis)
	registerCommonRoutes(router, nil, true, apis)
	registerUserRoutes(router, nil, true, apis)
	registerServerRoutes(router, nil, true)

	if apis.Admin == nil {
		t.Fatal("expected admin API to be registered")
	}

	if apis.Common == nil {
		t.Fatal("expected common API to be registered")
	}

	if len(apis.userAPIs) == 0 {
		t.Fatal("expected user sub-apis to be registered")
	}

	spec, err := apis.UserOpenAPI()
	if err != nil {
		t.Fatalf("expected merged user spec, got error: %v", err)
	}

	paths, ok := spec["paths"].(map[string]interface{})
	if !ok || len(paths) == 0 {
		t.Fatal("expected merged user spec to include paths")
	}

	if _, ok := paths["/api/v1/auth/login"]; !ok {
		t.Fatal("expected auth routes to be present in merged user spec")
	}

	if _, ok := paths["/api/v1/public/portal/purchase"]; !ok {
		t.Fatal("expected public portal routes to be present in merged user spec")
	}

	adminPaths := apis.Admin.OpenAPI().Paths
	if _, ok := adminPaths["/ads"]; ok {
		t.Fatal("expected ads routes to be removed from admin spec")
	}

	commonPaths := apis.Common.OpenAPI().Paths
	if _, ok := commonPaths["/ads"]; ok {
		t.Fatal("expected ads routes to be removed from common spec")
	}
}
