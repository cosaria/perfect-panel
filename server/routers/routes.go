package handler

import (
	"encoding/json"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gin-gonic/gin"
	"github.com/perfect-panel/server/routers/response"
	"github.com/perfect-panel/server/svc"
)

var bearerSecurity = []map[string][]string{{"bearer": {}}}

// apiConfig wraps apiConfig with $schema injection disabled.
// huma's default CreateHooks register a SchemaLinkTransformer that injects
// a "$schema" property into every response type — noise for SDK generation.
func apiConfig(title, version string) huma.Config {
	cfg := huma.DefaultConfig(title, version)
	cfg.CreateHooks = nil
	return cfg
}

func securitySchemes() map[string]*huma.SecurityScheme {
	return map[string]*huma.SecurityScheme{
		"bearer": {Type: "http", Scheme: "bearer", BearerFormat: "JWT"},
	}
}

// APIs holds all huma API instances for OpenAPI spec export.
type APIs struct {
	Admin    huma.API
	Common   huma.API
	userAPIs []huma.API // auth + public sub-APIs, merged via UserOpenAPI()
}

// UserOpenAPI merges all auth + public sub-API specs into a single OpenAPI spec.
func (a *APIs) UserOpenAPI() (map[string]interface{}, error) {
	merged := map[string]interface{}{
		"openapi": "3.1.0",
		"info":    map[string]interface{}{"title": "Perfect Panel User API", "version": "1.0.0"},
		"paths":   map[string]interface{}{},
		"components": map[string]interface{}{
			"schemas": map[string]interface{}{},
			"securitySchemes": map[string]interface{}{
				"bearer": map[string]interface{}{
					"type":         "http",
					"scheme":       "bearer",
					"bearerFormat": "JWT",
				},
			},
		},
	}

	paths := merged["paths"].(map[string]interface{})
	schemas := merged["components"].(map[string]interface{})["schemas"].(map[string]interface{})

	for _, api := range a.userAPIs {
		data, err := json.Marshal(api.OpenAPI())
		if err != nil {
			return nil, err
		}
		var spec map[string]interface{}
		if err := json.Unmarshal(data, &spec); err != nil {
			return nil, err
		}

		// Extract server prefix for absolute path construction
		prefix := ""
		if servers, ok := spec["servers"].([]interface{}); ok && len(servers) > 0 {
			if s, ok := servers[0].(map[string]interface{}); ok {
				prefix, _ = s["url"].(string)
			}
		}

		if specPaths, ok := spec["paths"].(map[string]interface{}); ok {
			for path, item := range specPaths {
				paths[prefix+path] = item
			}
		}

		if comps, ok := spec["components"].(map[string]interface{}); ok {
			if specSchemas, ok := comps["schemas"].(map[string]interface{}); ok {
				for name, schema := range specSchemas {
					schemas[name] = schema
				}
			}
		}
	}

	return merged, nil
}

func RegisterHandlers(router *gin.Engine, serverCtx *svc.ServiceContext) {
	registerHandlers(router, serverCtx, false)
}

// RegisterHandlersForSpec registers only route metadata (no middleware, no server routes).
// Used by the openapi export command — serverCtx can be nil.
func RegisterHandlersForSpec(router *gin.Engine) *APIs {
	return registerHandlers(router, nil, true)
}

func registerHandlers(router *gin.Engine, serverCtx *svc.ServiceContext, specOnly bool) *APIs {
	response.InstallHumaProblemFactory()

	apis := &APIs{}

	registerAdminRoutes(router, serverCtx, specOnly, apis)

	registerCommonRoutes(router, serverCtx, specOnly, apis)
	registerUserRoutes(router, serverCtx, specOnly, apis)
	registerServerRoutes(router, serverCtx, specOnly)

	return apis
}
