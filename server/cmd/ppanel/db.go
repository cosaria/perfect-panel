package ppanel

import (
	"fmt"

	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/internal/platform/support/conf"
	"github.com/perfect-panel/server/internal/platform/support/orm"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

func init() {
	dbCmd.PersistentFlags().StringVar(&dbConfigPath, "config", "etc/ppanel.yaml", "ppanel.yaml 文件路径")
	dbCmd.PersistentFlags().StringVar(&dbRevisionSource, "source", "embedded", "schema revision source")
}

var (
	dbConfigPath     string
	dbRevisionSource string
)

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "数据库 schema 工具",
	Long:  "管理 schema bootstrap、revisions、seed 和 reset 工作流。",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func openDBConnection(path string) (*gorm.DB, error) {
	var cfg config.File
	if err := conf.Load(path, &cfg); err != nil {
		return nil, fmt.Errorf("load config %s: %w", path, err)
	}

	return orm.ConnectMysql(orm.Mysql{
		Config: cfg.MySQL,
	})
}
