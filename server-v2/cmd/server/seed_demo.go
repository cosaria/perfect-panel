package main

import "github.com/perfect-panel/server-v2/internal/app/runtime"

func newSeedDemoCommand() subCommand {
	return subCommand{
		name:  runtime.ModeSeedDemo,
		short: "写入演示种子数据（骨架）",
		run: func() error {
			return nil
		},
	}
}
