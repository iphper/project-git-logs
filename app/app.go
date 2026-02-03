package app

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"project-git-logs/module"
	"project-git-logs/utils"
	"strings"
	"time"
)

// 日志消息类型
type LogMsg struct {
	name string // 项目名称
	log  string // 日志记录
}

type Application struct {
	user    string              // 获取提交用户信息
	dir     string              // 读取日志的目录
	date    []string            // 日期范围
	logChan chan LogMsg         // 日志记录通道
	history map[string]struct{} // 历史记录
	stop    chan struct{}       // 停止信号
}

// 初始化方法
func (app *Application) Init() {
	app.logChan = make(chan LogMsg, 10)
	app.stop = make(chan struct{})
	app.history = map[string]struct{}{}

	// 获取当前目录
	app.initDir()

	// 初始化git配置信息
	app.initGitConfig()

	// 初始化日期范围
	app.initDateRange()

}

// 初始化目录
func (app *Application) initDir() {
	app.dir = utils.Cmd("pwd")
}

// 初始化git配置
func (app *Application) initGitConfig() {
	// 获取git配置信息
	app.user = strings.ReplaceAll(utils.GitGlobalUserName(), "\n", "")
	if app.user == "" {
		app.user = module.NewInput("用户名：")
	}
	if app.user == "" {
		// 需要从终端获取
		panic("获取日志信息失败；读取不到提交用户信息")
	}
}

// 初始化日期范围
func (app *Application) initDateRange() {
	begin := utils.GetWeekStartDate()
	end := utils.GetTodayDate()
	// 默认日期范围
	defRange := begin + "~" + end

	dateRange := module.NewInput(fmt.Sprintf("日期范围[%s]：", defRange))
	if dateRange == "" {
		dateRange = defRange
	}

	// ~开头或者没有~的均为开始日期
	if !strings.Contains(dateRange, "~") || strings.HasPrefix(dateRange, "~") {
		begin = strings.TrimPrefix(dateRange, "~")
		// 日志日期范围
		app.date = []string{begin, end}
	} else if strings.HasSuffix(dateRange, "~") {
		end = strings.TrimSuffix(dateRange, "~")
		// 日志日期范围
		app.date = []string{begin, end}
	} else if strings.Contains(dateRange, "~") {
		app.date = strings.Split(dateRange, "~")[:2]
	}

	// 对日期进行格式化【只收日期，不收时间】
	for i, date := range app.date {
		if strings.Contains(date, " ") {
			app.date[i] = strings.Split(date, " ")[0]
		}
	}

}

// 日志读取
func (app *Application) worker(repo string, progress *module.Progress) {
	defer func() {
		wg.Done()
		progress.Add(1)
	}()
	repo = strings.TrimSpace(repo)
	repository := path.Base(repo)
	app.logChan <- LogMsg{
		name: repository,
		log: utils.PathCmd(
			repo,
			"git",
			"log",
			"--pretty=format:%s",
			"--no-merges",
			"--author", app.user,
			"--after", app.date[0]+" 00:00:00",
			"--before", app.date[1]+" 23:59:59",
		),
	}
}

// 日志写入
func (app *Application) writeWorker() {
	defer wg.Done()

	// 获取当前可执行文件所处目录
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}

	// 真实路径（处理符号链接）
	exePath, err = filepath.EvalSymlinks(exePath)
	if err != nil {
		panic(err)
	}

	exeDir := filepath.Dir(exePath)

	logFile := filepath.Join(exeDir, time.Now().Format("2006-01/02")+".log")

	// 判断目录是否存在
	dir := path.Dir(logFile)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, 0755)
	}

	fp, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}

	defer fp.Close()

	// 格式化写入日志文件
	fn := func(name, content string) {
		content = strings.TrimSpace(content)
		content = utils.GitFormat(content)
		content = strings.TrimSpace(content)

		if len(content) == 0 {
			return
		}

		// ;替换分号为换行符
		content = strings.ReplaceAll(content, ";", "\n")
		content = strings.ReplaceAll(content, "：", ":")

		for _, item := range strings.Split(content, "\n") {
			// 格式化
			item = strings.TrimSpace(item)
			if len(item) == 0 {
				continue
			}
			if _, ok := app.history[item]; ok {
				continue
			}
			app.history[item] = struct{}{}
			item = utils.GitFormat(item)
			fp.WriteString(fmt.Sprintf("[%v] %v\n", name, item))
		}
	}

	// 记录日志
	for {
		select {
		case <-app.stop: // 结束信号
			// 判断是否还有未写入的日志
			for len(app.logChan) > 0 {
				logMsg := <-app.logChan
				fn(logMsg.name, logMsg.log)
			}
			// 完成日志读取
			close(app.logChan)
			// 打开文件
			utils.Cmd("wsl.exe", "Code", logFile)
			return
		case logMsg, ok := <-app.logChan:
			if !ok {
				break
			}
			fn(logMsg.name, logMsg.log)
		}
	}
}

// Run 启动应用程序
func (app *Application) Run() {
	// waitgroup

	app.Init() // 初始化

	// 获取项目列表
	repos := utils.GitRepositories(strings.TrimSpace(app.dir))
	reposLen := len(repos)

	// 创建进度条
	progress := module.NewProgress().SetCount(float64(reposLen))

	wg.Add(2 + reposLen)
	go func() {
		defer wg.Done()
		progress.Run()
		app.stop <- struct{}{}
	}()

	// 启动一个协程记录日志
	go app.writeWorker()

	// 单独起协程处理一个仓库
	for _, rep := range repos {
		go app.worker(rep, progress)
	}

	// 等待协程结束
	wg.Wait()
}
