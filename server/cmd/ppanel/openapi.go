package ppanel

import (
	httpopenapi "github.com/perfect-panel/server/internal/platform/http/openapi"
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
		return httpopenapi.Export(outputDir)
	},
}
