package ppanel

import (
	"fmt"

	"github.com/perfect-panel/server/internal/platform/persistence/schema"
	schemarevisions "github.com/perfect-panel/server/internal/platform/persistence/schema/revisions"
	"github.com/spf13/cobra"
)

func init() {
	dbCmd.AddCommand(dbBootstrapCmd)
}

var dbBootstrapCmd = &cobra.Command{
	Use:   "bootstrap",
	Short: "初始化 schema registry 和 baseline revision",
	RunE: func(cmd *cobra.Command, args []string) error {
		schemarevisions.RegisterEmbedded()
		if err := schema.ValidateRevisionSource(dbRevisionSource); err != nil {
			return err
		}

		db, err := openDBConnection(dbConfigPath)
		if err != nil {
			return err
		}
		if err := schema.Bootstrap(db, dbRevisionSource); err != nil {
			return err
		}
		fmt.Println("[db bootstrap] schema bootstrap complete")
		return nil
	},
}
