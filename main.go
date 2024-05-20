package main

import (
	"fmt"
	"gitlog/libs"
	"path"

	"github.com/schollz/progressbar/v3"
)

// @func main
// @desc 入口函数
func main() {

	// 初始化
	initialize()

	// 项目列表
	projects := libs.GetProjectList()

	// 进度条
	bar := progressbar.Default(int64(len(projects)))

	// 创建所有日志文件
	libs.CreateAllLogFile()
	defer libs.CloseAllLogFile()

	// 日志前置内容写入
	libs.WriteLogBefor()
	// 遍历所有项目目录获取提交日志
	for project := range projects {
		bar.Describe(fmt.Sprintf("正在获取<%v>项目提交日志", path.Base(project)))
		libs.ReadProjectLog(project)
		bar.Add(1)
	}

	// 写入合并日志
	if libs.Merge {
		libs.WriteMergeLogs()
	}

	// 日志后置内容写入
	libs.WriteLogAfter()
	bar.Finish()

	// 后置操作
	libs.AfterCommand()
}

// @func initialize
// @desc 初始化函数
// @params void
// @returns void
func initialize() {
	// 初始化命令行参数
	libs.InitArgs()
	// 初始化配置
	libs.InitConfig()
}
