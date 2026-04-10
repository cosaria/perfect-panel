package main

import "github.com/perfect-panel/server-v2/internal/app/runtime"

func newMigrateCommand() subCommand {
	return subCommand{
		name:  runtime.ModeMigrate,
		short: "执行数据库迁移（骨架）",
		run: func() error {
			return nil
		},
	}
}
