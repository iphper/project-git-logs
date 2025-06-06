package cmd

import "sync"

var (
	command *Command
	once    sync.Once
)

// @func Single 获取单例的Command实例
func Single() *Command {
	once.Do(func() {
		command = new(Command)
		command.Init()
	})
	return command
}
