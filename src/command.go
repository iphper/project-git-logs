package src

import (
	"os"
	"os/exec"
)

// 运行全局命令
func command(cmdStr string, args ...string) string {
	fullpath, _ := os.Getwd()
	return pathCmd(fullpath, cmdStr, args...)
}

// 指定路径运行命令
func pathCmd(fullpath string, cmdStr string, args ...string) string {
	cmd := exec.Command(cmdStr, args...)
	cmd.Dir = fullpath
	res, _ := cmd.Output()
	if length := len(res); length > 0 {
		return string(res[:len(res)-1])
	}
	return ""
}
