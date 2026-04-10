package main

import "github.com/perfect-panel/server-v2/internal/app/runtime"

func newServeSchedulerCommand() subCommand {
	return subCommand{
		name:  runtime.ModeServeScheduler,
		short: "启动调度服务（骨架）",
		run: func() error {
			return nil
		},
	}
}
