package main

import "github.com/perfect-panel/server-v2/internal/app/runtime"

func newServeWorkerCommand() subCommand {
	return subCommand{
		name:  runtime.ModeServeWorker,
		short: "启动 Worker 服务（骨架）",
		run: func() error {
			return nil
		},
	}
}
