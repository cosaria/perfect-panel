package ppanel

import (
	"fmt"

	"github.com/perfect-panel/server/internal/platform/persistence/schema"
	schemarevisions "github.com/perfect-panel/server/internal/platform/persistence/schema/revisions"
	"github.com/spf13/cobra"
)

func init() {
	dbCmd.AddCommand(dbRevisionsCmd)
}

var dbRevisionsCmd = &cobra.Command{
	Use:   "revisions",
	Short: "执行 forward-only revisions",
	RunE: func(cmd *cobra.Command, args []string) error {
		schemarevisions.RegisterEmbedded()
		db, err := openDBConnection(dbConfigPath)
		if err != nil {
			return err
		}
		if err := schema.ApplyRevisions(db, dbRevisionSource); err != nil {
			return err
		}
		fmt.Println("[db revisions] revisions complete")
		return nil
	},
}
