package ppanel

import (
	"fmt"

	"github.com/perfect-panel/server/internal/platform/persistence/schema"
	schemarevisions "github.com/perfect-panel/server/internal/platform/persistence/schema/revisions"
	"github.com/perfect-panel/server/internal/platform/persistence/schema/seed"
	"github.com/spf13/cobra"
)

func init() {
	dbCmd.AddCommand(dbSeedCmd)
}

var dbSeedCmd = &cobra.Command{
	Use:   "seed",
	Short: "写入默认 seed 数据",
	RunE: func(cmd *cobra.Command, args []string) error {
		schemarevisions.RegisterEmbedded()
		db, err := openDBConnection(dbConfigPath)
		if err != nil {
			return err
		}
		if err := schema.Bootstrap(db, dbRevisionSource); err != nil {
			return err
		}
		if err := seed.Site(db); err != nil {
			return err
		}
		adminEmail, _ := cmd.Flags().GetString("admin-email")
		adminPassword, _ := cmd.Flags().GetString("admin-password")
		if err := seed.Admin(db, adminEmail, adminPassword); err != nil {
			return err
		}
		fmt.Println("[db seed] seed complete")
		return nil
	},
}

func init() {
	dbSeedCmd.Flags().String("admin-email", "admin@ppanel.dev", "管理员邮箱")
	dbSeedCmd.Flags().String("admin-password", "password", "管理员密码")
}
