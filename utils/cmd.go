package utils

import (
	"os/exec"
)

// 运行命令行命令
func Cmd(cmd string, args ...string) string {
	// 运行命令
	command := exec.Command(cmd, args...)
	output, err := command.CombinedOutput()

	// 执行命令
	if err != nil {
		return ""
	}
	return string(output)
}

func PathCmd(path string, cmd string, args ...string) string {
	// 切换到指定路径并运行命令
	command := exec.Command(cmd, args...)
	command.Dir = path // 设置工作目录
	output, err := command.CombinedOutput()

	// 执行命令
	if err != nil {
		return ""
	}
	return string(output)
}
