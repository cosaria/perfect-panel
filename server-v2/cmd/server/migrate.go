package main

import (
	"context"

	"github.com/perfect-panel/server-v2/internal/app/runtime"
	ppdb "github.com/perfect-panel/server-v2/internal/platform/db"
)

func newMigrateCommand() subCommand {
	return subCommand{
		name:  runtime.ModeMigrate,
		short: "执行数据库迁移（骨架）",
		run: func() error {
			database, err := ppdb.OpenFromEnv()
			if err != nil {
				return err
			}
			defer database.Close()

			return ppdb.Migrate(context.Background(), database)
		},
	}
}
