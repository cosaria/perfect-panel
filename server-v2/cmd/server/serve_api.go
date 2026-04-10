package main

import "github.com/perfect-panel/server-v2/internal/app/runtime"

func newServeAPICommand() subCommand {
	return subCommand{
		name:  runtime.ModeServeAPI,
		short: "启动 API 服务（骨架）",
		run: func() error {
			return nil
		},
	}
}
