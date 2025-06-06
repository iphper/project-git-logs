package cmd

import (
	"os"
	"os/exec"
	"strings"
)

type Command struct {
	options map[string]any
	args    []any
}

// @func Init 初始化
func (c *Command) Init() {
	// 初始化
	c.options = map[string]any{}
	c.args = []any{}

	before := ""
	for _, val := range os.Args[1:] {

		// 如果是选项
		if strings.HasPrefix(val, "-") {
			// 去掉前缀
			val = strings.TrimLeft(val, "-")

			// 如果是选项且有值
			if strings.Contains(val, "=") {
				// 分割选项和值
				parts := strings.SplitN(val, "=", 2)
				if len(parts) == 2 {
					c.SetOption(parts[0], parts[1])
				} else {
					c.SetOption(parts[0], true)
				}
			} else {
				// 如果是选项但没有值
				before = val
			}
		} else if before != "" {
			c.SetOption(before, val)
			before = "" // 重置 before
		} else {
			c.AddArg(val)
		}
	}
}

// ======< command funcs >======

// @func Run 运行命令
func (c *Command) Run(cmd string, args ...string) string {
	// 运行命令
	command := exec.Command(cmd, args...)
	output, err := command.CombinedOutput()

	// 执行命令
	if err != nil {
		panic("执行命令失败: " + err.Error())
	}
	return string(output)
}

// @func PathRun 切换到指定路径并运行命令
func (c *Command) PathRun(path string, cmd string, args ...string) string {
	// 切换到指定路径并运行命令
	command := exec.Command(cmd, args...)
	command.Dir = path // 设置工作目录
	output, err := command.CombinedOutput()

	// 执行命令
	if err != nil {
		panic("执行命令失败: " + err.Error())
	}
	return string(output)
}

// ======< options funcs >======

// @func SetOption 设置选项
func (c *Command) SetOption(key string, value any) *Command {
	// 设置选项
	c.options[key] = value
	return c
}

// @func GetOption 获取选项
func (c *Command) GetOption(key string) any {
	// 获取选项
	if value, exists := c.options[key]; exists {
		return value
	}
	return nil
}

// ======< agrs funcs >======

// @func AddArg 添加参数
func (c *Command) AddArg(arg any) *Command {
	// 添加参数
	c.args = append(c.args, arg)
	return c
}

// @func GetArg 获取参数
func (c *Command) GetArg(idx int) any {
	// 获取参数
	return c.args[idx]
}
