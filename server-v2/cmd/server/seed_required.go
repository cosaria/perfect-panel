package main

import "github.com/perfect-panel/server-v2/internal/app/runtime"

func newSeedRequiredCommand() subCommand {
	return subCommand{
		name:  runtime.ModeSeedRequired,
		short: "写入必需种子数据（骨架）",
		run: func() error {
			return nil
		},
	}
}
