package main

import (
	"context"

	"github.com/perfect-panel/server-v2/internal/app/runtime"
	ppdb "github.com/perfect-panel/server-v2/internal/platform/db"
	"github.com/perfect-panel/server-v2/internal/platform/db/seeds"
)

func newSeedDemoCommand() subCommand {
	return subCommand{
		name:  runtime.ModeSeedDemo,
		short: "写入演示种子数据（骨架）",
		run: func() error {
			database, err := ppdb.OpenFromEnv()
			if err != nil {
				return err
			}
			defer database.Close()

			return seeds.ApplyDemo(context.Background(), database)
		},
	}
}
