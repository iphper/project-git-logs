package app

import (
	"fmt"
	"os"
	"path/filepath"
	"project-git-logs/module"
	"project-git-logs/utils"
	"runtime"
	"strings"
	"time"
)

// 日志消息类型
type LogMsg struct {
	name string // 项目名称
	log  string // 日志记录
}

// Application 应用
type Application struct {
	gitUser  string              // 获取提交用户信息
	root     string              // 读取日志的目录
	binRoot  string              // 可执行文件目录
	date     []string            // 日期范围
	logChan  chan LogMsg         // 日志记录通道
	history  map[string]struct{} // 历史记录
	process  int                 // 进程数据
	progress *module.Progress    // 进度条
	logFile  string              // 日志文件路径
}

// 应用执行方法
func (app *Application) Run() {

	defer func() {
		if a := recover(); a != nil {
			r := []byte{}
			n := runtime.Stack(r, true)
			fmt.Printf("panic: %v\nstack: %s\n", a, string(r[:n]))
			fmt.Println("发生错误：", a)
			os.Exit(1)
		}
	}()

	app.Init()

	app.Worker()

	app.AfterWorker()
}

// 初始化
func (app *Application) Init() {
	// 初始化日志通道和历史记录
	app.logChan = make(chan LogMsg, 10)
	app.history = map[string]struct{}{}

	// 获取当前目录
	app.root, _ = os.Getwd()
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
	// 可执行文件目录
	app.binRoot = filepath.Dir(exePath)

	// 初始化git配置信息
	app.gitUser = strings.ReplaceAll(utils.GitGlobalUserName(), "\n", "")
	if app.gitUser == "" {
		app.gitUser = module.NewInput("用户名：")
	}
	if app.gitUser == "" {
		// 需要从终端获取
		panic("获取日志信息失败；读取不到提交用户信息")
	}

	begin := utils.GetWeekStartDate()
	end := utils.GetTodayDate()
	// 默认日期范围
	defRange := begin + "~" + end

	dateRange := ""

	// 默认为第一个参数
	if l := len(os.Args); l > 1 {
		dateRange = os.Args[1]
	} else {
		dateRange = module.NewInput(fmt.Sprintf("日期范围[%s]：", defRange))
	}
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

// 读取
func (app *Application) Read(repo string) {
	defer func() {
		app.progress.Add(1)
		app.process-- // 完成后更新计数
		wg.Done()
	}()
	repo = strings.TrimSpace(repo)
	// 适配win环境下的路径【path.Base在win下会将绝对路径当作仓库名称】
	repo = strings.ReplaceAll(repo, "\\", "/")
	repository := filepath.Base(repo)
	app.logChan <- LogMsg{
		name: repository,
		log: utils.PathCmd(
			repo,
			"git",
			"log",
			"--pretty=format:%s",
			"--no-merges",
			"--author", app.gitUser,
			"--after", app.date[0]+" 00:00:00",
			"--before", app.date[1]+" 23:59:59",
		),
	}
}

// 写入
func (app *Application) Write() {
	defer wg.Done()

	app.logFile = filepath.Join(app.binRoot, time.Now().Format("2006-01/02")+".log")

	// 判断目录是否存在
	dir := filepath.Dir(app.logFile)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, 0755)
	}

	fp, err := os.OpenFile(app.logFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
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

	// 退出函数处理
	exitFn := func() {
		// 判断是否还有未写入的日志
		for len(app.logChan) > 0 {
			logMsg := <-app.logChan
			fn(logMsg.name, logMsg.log)
		}
		// 完成日志读取
		close(app.logChan)
	}

	// 记录日志
	for {
		if app.process <= 0 {
			exitFn()
			break
		}
		select {
		case logMsg, ok := <-app.logChan:
			if !ok {
				break
			}
			fn(logMsg.name, logMsg.log)
		case <-time.After(time.Millisecond * 50): // 超时跳过防止等待太久
			// fmt.Println("等待...")
		}
	}
}

// Worker
func (app *Application) Worker() {

	// 读取所有git仓库
	repos := utils.GitRepositories(strings.TrimSpace(app.root))
	app.process = len(repos)

	// 开启进度条
	app.progress = module.NewProgress(float64(app.process))
	wg.Add(2 + app.process) // 2个协程 + 每个仓库一个协程
	go func() {
		defer wg.Done()
		app.progress.Run()
	}()

	// 开启写入任务
	go app.Write()

	// 开始读取任务
	for _, repo := range repos {
		go app.Read(repo)
	}

	wg.Wait()
}

// Worker结束后的处理
func (app *Application) AfterWorker() {
	// 如果是windows系统
	if runtime.GOOS == "windows" {
		// 以非阻塞方式打开文件
		utils.SyncCmd("Code", app.logFile)
	} else {
		// 以非阻塞方式通过 wsl 启动
		utils.SyncCmd("wsl.exe", "Code", app.logFile)
	}
}
