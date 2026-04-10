package main

import (
	"fmt"
	"io"
	"slices"
	"strings"
)

type subCommand struct {
	name  string
	short string
	run   func() error
}

type rootCommand struct {
	use         string
	short       string
	subCommands []subCommand
}

func Execute(args []string, stdout io.Writer, stderr io.Writer) error {
	_ = stderr

	root := newRootCommand()
	if len(args) == 0 {
		root.printHelp(stdout)
		return nil
	}
	if isHelpToken(args[0]) {
		if len(args) > 1 {
			return fmt.Errorf("根命令 help 不接受额外参数: %s", strings.Join(args[1:], " "))
		}
		root.printHelp(stdout)
		return nil
	}

	name := args[0]
	for _, cmd := range root.subCommands {
		if cmd.name != name {
			continue
		}
		if len(args) > 1 && isHelpToken(args[1]) {
			if len(args) > 2 {
				return fmt.Errorf("命令 %s 不接受额外参数: %s", cmd.name, strings.Join(args[2:], " "))
			}
			root.printSubCommandHelp(stdout, cmd)
			return nil
		}
		if len(args) > 1 {
			return fmt.Errorf("命令 %s 不接受额外参数: %s", cmd.name, strings.Join(args[1:], " "))
		}
		return cmd.run()
	}

	root.printHelp(stdout)
	return fmt.Errorf("未知命令: %s", name)
}

func newRootCommand() rootCommand {
	return rootCommand{
		use:   "server",
		short: "Perfect Panel v2 服务入口",
		subCommands: []subCommand{
			newServeAPICommand(),
			newServeWorkerCommand(),
			newServeSchedulerCommand(),
			newMigrateCommand(),
			newSeedRequiredCommand(),
			newSeedDemoCommand(),
		},
	}
}

func (c rootCommand) printHelp(stdout io.Writer) {
	fmt.Fprintf(stdout, "%s\n\n", c.short)
	fmt.Fprintf(stdout, "用法:\n  %s <command>\n\n", c.use)
	fmt.Fprintln(stdout, "可用命令:")
	for _, cmd := range c.subCommands {
		fmt.Fprintf(stdout, "  %-16s %s\n", cmd.name, cmd.short)
	}
	fmt.Fprintln(stdout)
	fmt.Fprintln(stdout, "使用 \"server <command> --help\" 查看子命令帮助。")
}

func (c rootCommand) printSubCommandHelp(stdout io.Writer, cmd subCommand) {
	fmt.Fprintf(stdout, "%s\n\n", cmd.short)
	fmt.Fprintf(stdout, "用法:\n  %s %s\n", c.use, cmd.name)
}

func isHelpToken(token string) bool {
	return slices.Contains([]string{"--help", "-h", "help"}, strings.TrimSpace(token))
}
