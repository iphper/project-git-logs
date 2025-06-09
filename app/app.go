package app

import (
	"fmt"
	"path"
	"project-git-logs/cmd"
	"project-git-logs/config"
	"project-git-logs/git"
	"project-git-logs/utils"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
)

type App struct {
	cmd    *cmd.Command
	config *config.Config // 配置
}

// @func Init 初始化
func (a *App) Init() {
	// 初始化命令行参数
	a.cmd = cmd.Single()
	// 初始化配置
	a.config = config.Single()
}

// @func Error 错误处理
func (a *App) Error(err error) {
	// 错误处理
	fmt.Printf("发生错误: %s\n", err.Error())
}

// @func Start 启动

// @func Start 启动
func (a *App) Start() {
	// 启动前的准备工作
	before := a.cmd.GetOption("before")
	date := utils.GetWeekStartDate() + " 00:00:00"
	if before == nil {
		before = utils.Scan("请输入提交日期范围的起始日期["+date+"]: ", date)
		if !utils.IsValidDate(before.(string)) {
			before = date
		}
		a.cmd.SetOption("before", before)
	}

	// 获取提交日期范围
	after := a.cmd.GetOption("after")
	if after == nil {
		date = utils.GetTodayDate() + " 23:59:59"
		after = utils.Scan("请输入提交日期范围的结束日期["+date+"]: ", date)
		if !utils.IsValidDate(after.(string)) {
			after = date
		}
		a.cmd.SetOption("after", after)
	}

	// 提交账号
	if a.cmd.GetOption("author") == nil {
		a.cmd.SetOption("author", a.config.User)
	}

	// 读取所有git仓库
	dirs := git.ReadGitPaths(a.config.Path)

	wg := sync.WaitGroup{}

	length := len(dirs)
	// 协程组数
	wg.Add(length)
	// 进度条
	bar := progressbar.NewOptions(length, progressbar.OptionFullWidth())

	for _, dir := range dirs {
		go func(dir string) {
			// 协程完成
			defer wg.Done()
			// 进度条改变
			defer bar.Add(1)

			// 读取git仓库提交日志
			if err := git.ReadGitLogs(dir); err != nil {
				a.Error(err)
			}
		}(dir)
	}
	// 协程等待
	wg.Wait()
	// 进度条完成
	bar.Finish()

	// 写入所有git仓库的提交日志
	logFile := path.Join(a.config.LogPath, time.Now().Format("2006-01/02.log"))
	git.WriteGitLogs(logFile)

	// 如果打开文件
	if a.config.IsOpen() {
		cmd.Single().Run("code", logFile)
	}
}

// @func Run 运行
func (a *App) Run() {

	// 错误处理
	defer func() {
		if err := recover(); err != nil {
			a.Error(err.(error))
		}
	}()

	// 初始化
	a.Init()

	// 启动
	a.Start()
}
