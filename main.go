package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/schollz/progressbar/v3"
)

// 常量
const (
	// 获取配置文件【保存一些配置到此文件】
	// 如果存在此文件就获取此文件中的参数
	// 如果不存在初始化时需要使用者输入
	configFileName = "weeker.conf"
)

// 变量
var (
	// 配置数据
	config map[string][]string = make(map[string][]string, 1)
	// 日志项编号
	logIdenxNO = 1
	// 开始时间
	after = time.Now().AddDate(0, 0, -int(time.Now().Weekday()-time.Monday)).Format("2006-01-02") + " 00:00:00"
	// 结束时间
	before = time.Now().Format("2006-01-02") + " 23:59:59"
	// 日志文件句柄列表
	logFiles []*os.File
	// 日志文件列表
	logPaths []string
)

// @func main
// @desc 入口函数
func main() {
	// 初始化
	initialize()

	// 项目列表
	projects := getProjectList()

	// 进度条
	bar := progressbar.Default(int64(len(projects)))

	// 创建所有日志文件
	createAllLogFile()
	defer closeAllLogFile()

	// 日志前置内容写入
	writeLogBefor()
	// 遍历所有项目目录获取提交日志
	for project := range projects {
		bar.Describe(fmt.Sprintf("正在获取<%v>项目提交日志", path.Base(project)))
		readProjectLog(project)
		bar.Add(1)
	}
	// 日志后置内容写入
	writeLogAfter()
	bar.Finish()

	// 后置操作
	afterCommand()
}

// @func loadProjects
// @desc 加载项目目录
// @param projectPath string
// @return []string
func loadProjects(projectPath string) (list []string) {
	if isProject(projectPath) {
		list = append(list, projectPath)
	} else {
		dirs, _ := os.ReadDir(projectPath)
		var temp string

		for _, project := range dirs {
			temp = path.Join(projectPath, project.Name())
			if isProject(temp) {
				list = append(list, temp)
				continue
			}
			if project.IsDir() {
				list = append(list, loadProjects(temp)...)
			}
		}
	}
	return
}

// @func isProject
// @desc 是否为项目目录
// @param name string
// @return bool
func isProject(name string) bool {
	dirs, _ := os.ReadDir(name)
	str := ".git,README.md,.idea,"
	for _, item := range dirs {
		if strings.Contains(str, item.Name()+",") {
			return true
		}
	}
	return false
}

// @func initialize
// @desc 初始化函数
// @params void
// @returns void
func initialize() {
	var input string
	// 获取获取日期区间
	fmt.Printf("请输入读取提交日志的开始时间[%v]:", after)
	fmt.Scanln(&input)
	if input != "" {
		after = input
	}
	// 获取获取日期区间
	fmt.Printf("请输入读取提交日志的结束时间[%v]:", before)
	fmt.Scanln(&input)
	if input != "" {
		before = input
	}

	// 配置参数获取
	if !pathExists(configFileName) {
		initConfig()
	} else {
		readConfig()
	}

	// 配置占位符转化
	for field, values := range config {
		for idx, value := range values {
			switch true {
			// 时间日期点位符转化
			case strings.Contains(value, "$DATE"):
				config[field][idx] = strings.ReplaceAll(value, "$DATE", time.Now().Format("2006-01-02"))
			case strings.Contains(value, "$TIME"):
				config[field][idx] = strings.ReplaceAll(value, "$DATE", time.Now().Format("15:04:05"))
			case strings.Contains(value, "$DATETIME"):
				config[field][idx] = strings.ReplaceAll(value, "$DATE", time.Now().Format("2006-01-02 15:04:05"))
			}
		}
	}
}

// @func initConfig
// @desc 初始化配置
// @param void
// @return void
func initConfig() {
	fmt.Println("-----------------------------")
	fmt.Println("未找到配置信息，请填写以下配置项")
	keys := map[string]string{
		"project_root": "项目路径",
		"author":       "提交用户名",
		"log_root":     "日志文件保存路径",
		"filename":     "文件名格式",
		"log_before":   "日志前置内容",
		"log_after":    "日志后置内容",
		"after_cmd":    "完成后执行命令操作",
	}

	var input string
	var values []string
	var i int

	for field, desc := range keys {
		values = []string{}
		for i = 1; ; {
			input = ""
			fmt.Printf("请输入第%d个%v[回车确认]:", i, desc)
			fmt.Scanln(&input)
			if input == "" {
				break
			}
			values = append(values, input)
			i++
		}

		config[field] = values[:i-1]
	}

	// 询问是否写入配置文件
	for {
		fmt.Print("是否保存到配置文件？[yes|no]:")
		fmt.Scanln(&input)

		input = strings.ToLower(input)
		if input == "yes" || input == "y" {
			// 保存到配置文件
			writeConfig()
			break
		}
		if input == "no" || input == "n" {
			// 不保存
			break
		}
	}
	fmt.Println("-----------------------------")
}

// @func readConfig
// @desc 读取配置文件
// @param void
// @return void
func readConfig() {
	var err error
	// 存在配置文件就读取配置文件内容
	bufs, e := os.ReadFile(configFileName)
	if e != nil {
		handleError(err)
	}

	// 读取配置
	err = json.Unmarshal(bufs, &config)
	if err != nil {
		handleError(err)
	}
}

// @func writeConfig
// @desc 写入配置文件
// @param void
// @return void
func writeConfig() {
	file, err := os.Create(configFileName)
	if err != nil {
		handleError(err)
	}

	defer file.Close()

	bufs, err := json.Marshal(config)
	if err != nil {
		handleError(err)
	}

	file.Write(bufs)
}

// @func pathExists
// @desc 路径是否存在
// @param path string
// @return bool
func pathExists(path string) bool {
	if len(path) < 1 {
		return false
	}

	_, e := os.Stat(path)
	if e == nil {
		return true
	}
	return os.IsExist(e)
}

// @func handleError
// @desc 处理错误
// @param err error
// @return void
func handleError(err error) {
	fmt.Println(err.Error())
	os.Exit(1)
}

// @func getProjectList
// @desc 获取项目列表
// @param void
// @return []string
func getProjectList() map[string]uint {

	// 获取所有的项目目录列表
	projects := map[string]uint{}
	for _, item := range config["project_root"] {
		for _, project := range loadProjects(item) {
			projects[project] = 1
		}
	}

	return projects
}

// @func createAllLogFile
// @desc 创建所有日志文件句柄
// @param void
// @return void
func createAllLogFile() {
	filename := strings.Join(config["filename"], "") + ".log"
	for _, dir := range config["log_root"] {
		dir = path.Join(dir, time.Now().Format("2006-01"))
		// 没有目录则创建
		if !pathExists(dir) {
			os.MkdirAll(dir, os.ModePerm)
		}
		path := path.Join(dir, filename)
		logPaths = append(logPaths, path)
		file, err := os.Create(path)
		if err != nil {
			// 创建错误时跳过
			fmt.Println(err)
			continue
		}
		logFiles = append(logFiles, file)
	}
}

// @func closeAllLogFile
// @desc 关闭所有文件句柄
// @param void
// @return void
func closeAllLogFile() {
	for _, file := range logFiles {
		file.Close()
	}
}

// @func writeLogBefor
// @desc 写入日志前置内容
// @param void
// @return void
func writeLogBefor() {
	for _, item := range config["log_before"] {
		writeLog(fmt.Sprintf("%s\n", item))
	}
}

// @func writeLogAfter
// @desc 写入日志后置内容
// @param void
// @return void
func writeLogAfter() {
	for _, item := range config["log_after"] {
		writeLog(fmt.Sprintf("%s\n", item))
	}
}

// @func readProjectLog
// @desc 读取项目日志
// @param project string
// @returns void
func readProjectLog(project string) {
	name := path.Base(project)
	// 日志
	gitLogArgs := []string{
		"log",
		"--oneline",
		"--no-merges",
		"--all",
		"--after",
		after,
		"--before",
		before,
	}

	// 指定author
	for _, author := range config["author"] {
		gitLogArgs = append(gitLogArgs, "--author", author)
	}

	cmd := exec.Command("git", gitLogArgs...)
	cmd.Dir = project

	temp, _ := cmd.Output()
	logs := strings.Split(string(temp), "\n")
	temp = nil

	// 倒序
	for i, j := 0, len(logs)-1; i < j; i, j = i+1, j-1 {
		logs[i], logs[j] = logs[j], logs[i]
	}

	// 遍历
	for _, line := range logs {
		line = strings.Replace(line, "feat:", "", 1)
		for _, sp := range []string{" ", "；"} {
			line = strings.ReplaceAll(line, sp, ";")
		}
		items := strings.Split(line, ";")[1:]
		for _, item := range items {
			if len(item) == 0 {
				continue
			}
			writeLog(fmt.Sprintf("%3d、[%v] %v\n", logIdenxNO, name, item))
			logIdenxNO++
		}
	}
}

// @func writeLog
// @desc 写入日志
// @param name string 项目名称
// @param log string 日志文本
// @returns void
func writeLog(log string) {
	for _, file := range logFiles {
		file.WriteString(log)
	}
}

// @func afterCommand
// @desc 后置操作
// @param void
// @return void
func afterCommand() {
	for _, cmd := range config["after_cmd"] {
		for _, path := range logPaths {
			exec.Command(cmd, path).Run()
		}
	}
}
