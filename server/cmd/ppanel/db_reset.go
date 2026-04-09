package ppanel

import (
	"fmt"

	"github.com/perfect-panel/server/internal/platform/persistence/schema"
	schemarevisions "github.com/perfect-panel/server/internal/platform/persistence/schema/revisions"
	"github.com/spf13/cobra"
)

func init() {
	dbCmd.AddCommand(dbResetCmd)
}

var dbResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "重置 schema registry 和 baseline revision",
	RunE: func(cmd *cobra.Command, args []string) error {
		schemarevisions.RegisterEmbedded()
		db, err := openDBConnection(dbConfigPath)
		if err != nil {
			return err
		}
		if err := schema.Reset(db, dbRevisionSource); err != nil {
			return err
		}
		fmt.Println("[db reset] schema reset complete")
		return nil
	},
}
