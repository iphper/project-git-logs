package src

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/schollz/progressbar/v3"
)

type gitLog struct {
	// 配置文件
	conf string
	// 缓存文件
	config map[string][]string
	// 日志文件
	logfile *os.File
	// 文件名称
	filename string
	// 项目列表
	list []string
	// 进度条
	bar *progressbar.ProgressBar
}

// 流程
// 1、查询配置文件
// 1.1、有则跳到下一步；没有则创建
// 1.2、配置信息[项目目录、前置操作命令、后置操作命令、提交日志分隔符]
// 1.3、定义占位符及替换
// 1.4、获取git全局配置的用户名【用于获取提交的作者过滤】
func (app *gitLog) initialize() {

	// 读取配置文件
	if !PathExists(app.conf) {
		// 走创建
		cache.Create(app.conf)
	}

	// 读取配置
	app.config = cache.Read(app.conf)

	// 读取所有项目目录
	go func(app *gitLog) {
		for _, p := range app.config["root"] {
			app.list = append(app.list, LoadProjects(p)...)
		}
	}(app)

	// 读取git全局配置
	app.config["logcmd"] = []string{
		"log",
		"--oneline",
		"--no-merges",
		"--all",
		"--author",
		command("git", "config", "user.name"),
	}

}

// 2、前置操作命令执行
func (app *gitLog) before() {
	var err error
	// 日志文件名
	app.filename = path.Join(time.Now().Format("2006-01/02.log"))

	if !PathExists(path.Dir(app.filename)) {
		os.Mkdir(path.Dir(app.filename), 0777)
	}

	// 生成日志文件
	if app.logfile, err = os.Create(app.filename); err != nil {
		panic(err)
	}
}

// 3、输入获取的日志日期范围[默认本周一至周五]
// 3.1、获取所有项目目录
// 3.1、每获获取一个项目目录均开启一个协程读取日志
// 3.2、开启一个协程写入日志[日志文件为当天日期]
func (app *gitLog) run() {
	var err error

	// 默认日期范围
	def := []string{
		time.Now().AddDate(0, 0, -int(time.Now().Weekday()-time.Monday)).Format("2006-01-02"),
		time.Now().Format("2006-01-02"),
	}

	// 开始日期[周一]
	input := ""
	// 获取获取开始日期
	for {
		fmt.Printf("请输入读取提交日志的开始时间[%v]:", def[0])
		fmt.Scanln(&input)
		if input != "" {
			if _, err = time.Parse("2006-01-02", input); err != nil {
				app.config["logcmd"] = append(app.config["logcmd"], "--after", "'"+input[:10]+" 00:00:00'")
				break
			}
			fmt.Println("输入的日期格式不正确")
		} else {
			app.config["logcmd"] = append(app.config["logcmd"], "--after", "'"+def[0]+" 00:00:00'")
			break
		}
	}

	// 结束日期[当天]
	input = ""
	// 获取获取结束日期
	for {
		fmt.Printf("请输入读取提交日志的结束时间[%v]:", def[1])
		fmt.Scanln(&input)
		if input != "" {
			if _, err = time.Parse("2006-01-02", input); err != nil {
				app.config["logcmd"] = append(app.config["logcmd"], "--before", "'"+input[:10]+" 23:59:59'")
				break
			}
			fmt.Println("输入的日期格式不正确")
		} else {
			app.config["logcmd"] = append(app.config["logcmd"], "--before", "'"+def[1]+" 23:59:59'")
			break
		}
	}

	// 总进度
	app.bar = progressbar.Default(int64(len(app.list)))

	// 开启写入协程
	go app.writeLog()

	// 读取
	for _, log := range app.list {
		group.Add(1)
		go app.readLog(log)
	}

	// // 等待
	group.Wait()

	// 写入空chan
	writeChan <- []string{}

	// 完成
	app.bar.Finish()

}

// 4、后置操作命令执行
func (app *gitLog) after() {
	// 关闭文件
	app.logfile.Close()

	// 后置操作
	if len(app.config["after"]) > 0 {
		for _, cmd := range app.config["after"] {
			command(cmd, app.filename)
		}
	}

}

// 写入log文件
func (app *gitLog) writeLog() {
	defer app.logfile.Close()
	lineNO := 1
	line := <-writeChan
	for len(line) > 0 {
		app.logfile.WriteString(fmt.Sprintf("[%s] %d、%s\n", line[0], lineNO, line[1]))
		line = <-writeChan
		lineNO += 1
	}
	close(writeChan)
}

// 读取log文件
func (app *gitLog) readLog(logpath string) {
	defer group.Done()
	defer app.bar.Add(1)

	app.bar.Describe(fmt.Sprintf("正在获取<%v>项目提交日志", logpath))

	// 获取提交日志
	log := pathCmd(logpath, "git", app.config["logcmd"]...)

	log = strings.ReplaceAll(log, "\n", ";")
	log = strings.ReplaceAll(log, "；", ";")
	logs := strings.Split(log, ";")

	// 倒序
	for i, j := 0, len(logs)-1; i < j; i, j = i+1, j-1 {
		logs[i], logs[j] = logs[j], logs[i]
	}

	for _, line := range logs {
		// 跳过空行
		if len(line) <= 0 {
			continue
		}

		lineArr := strings.Split(line, " ")

		writeChan <- []string{path.Base(logpath), lineArr[len(lineArr)-1]}
	}

}

// 初始化
func init() {
	once.Do(func() {
		app = &gitLog{
			conf:   "config.data",
			config: map[string][]string{},
		}

		// 初始化
		app.initialize()
	})
}

// 执行操作
func Runing() {
	app.before()

	app.run()

	app.after()
}
