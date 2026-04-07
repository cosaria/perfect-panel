package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/perfect-panel/server/routers"
	"github.com/spf13/cobra"
)

func init() {
	openapiCmd.Flags().StringP("output", "o", "docs/openapi", "Output directory for spec files")
	rootCmd.AddCommand(openapiCmd)
}

var openapiCmd = &cobra.Command{
	Use:   "openapi",
	Short: "Export OpenAPI 3.1 specs",
	Long:  "Export OpenAPI 3.1 JSON specs from huma route definitions. No database connection required.",
	RunE: func(cmd *cobra.Command, args []string) error {
		outputDir, _ := cmd.Flags().GetString("output")
		if outputDir == "" {
			outputDir = "docs/openapi"
		}

		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("create output dir: %w", err)
		}

		gin.SetMode(gin.ReleaseMode)
		router := gin.New()

		// RegisterHandlersForSpec only registers route metadata (no middleware).
		// Handler functions are registered but never called during spec export.
		apis := handler.RegisterHandlersForSpec(router)

		userSpec, err := apis.UserOpenAPI()
		if err != nil {
			return fmt.Errorf("merge user spec: %w", err)
		}

		specs := map[string]interface{}{
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
			if err := os.WriteFile(path, data, 0644); err != nil {
				return fmt.Errorf("write %s: %w", path, err)
			}
			fmt.Printf("Exported %s (%d bytes)\n", path, len(data))
		}

		return nil
	},
}
