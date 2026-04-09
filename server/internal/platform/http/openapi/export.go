package openapi

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	handler "github.com/perfect-panel/server/internal/platform/http"
)

func Export(outputDir string) error {
	if outputDir == "" {
		outputDir = "docs/openapi"
	}

	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	apis := handler.RegisterHandlersForSpec(router)

	userSpec, err := apis.UserOpenAPI()
	if err != nil {
		return fmt.Errorf("merge user spec: %w", err)
	}

	specs := map[string]any{
		"admin":  apis.Admin.OpenAPI(),
		"common": apis.Common.OpenAPI(),
		"user":   userSpec,
	}

	for name, spec := range specs {
		data, err := json.MarshalIndent(spec, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal %s spec: %w", name, err)
		}

		path := filepath.Join(outputDir, name+".json")
		if err := os.WriteFile(path, data, 0o644); err != nil {
			return fmt.Errorf("write %s: %w", path, err)
		}
	}

	return nil
}
