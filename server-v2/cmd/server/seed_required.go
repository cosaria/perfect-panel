package main

import (
	"context"

	"github.com/perfect-panel/server-v2/internal/app/runtime"
	ppdb "github.com/perfect-panel/server-v2/internal/platform/db"
	"github.com/perfect-panel/server-v2/internal/platform/db/seeds"
)

func newSeedRequiredCommand() subCommand {
	return subCommand{
		name:  runtime.ModeSeedRequired,
		short: "写入必需种子数据（骨架）",
		run: func() error {
			database, err := ppdb.OpenFromEnv()
			if err != nil {
				return err
			}
			defer database.Close()

			return seeds.ApplyRequired(context.Background(), database)
		},
	}
}
