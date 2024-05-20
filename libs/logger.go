package libs

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"time"
)

var (
	logIdenxNO = 1
	// 日志文件句柄列表
	logFiles []*os.File
	// 日志文件列表
	logPaths []string
	// 开始时间
	after = DateFormat(time.Now().AddDate(0, 0, -int(time.Now().Weekday()-time.Monday)))
	// 结束时间
	before = DateFormat()
	// 合并日志
	mergeLogs map[string][]string = make(map[string][]string, 1)
	// 合并规则
	rulesMap = map[string][]func(string, string) (bool, string){
		"xcx": {
			func(s1 string, s2 string) (bool, string) {
				list := regexp.MustCompile(`《\S+》`).FindAllString(s2, -1)
				if len(list) > 0 {
					return true, strings.ReplaceAll(strings.ReplaceAll(list[0], `《`, ""), `》`, "")
				}
				return false, ""
			},
		},
	}
	// 合并写入
	writeMap = map[string]func(string, []string){
		"xcx": func(s1 string, s2 []string) {
			log := fmt.Sprintf("%3d、[%v] %v\n", logIdenxNO, "小程序", strings.Join(s2, ","))
			WriteLog(log)
			logIdenxNO++
		},
	}
)

// @func createAllLogFile
// @desc 创建所有日志文件句柄
// @param void
// @return void
func CreateAllLogFile() {
	filename := strings.Join(Config["filename"], "") + ".txt"
	for _, dir := range Config["log_root"] {
		dir = path.Join(dir, time.Now().Format("2006-01"))
		// 没有目录则创建
		if !PathExists(dir) {
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
func CloseAllLogFile() {
	for _, file := range logFiles {
		file.Close()
	}
}

// @func writeLogBefor
// @desc 写入日志前置内容
// @param void
// @return void
func WriteLogBefor() {
	for _, item := range Config["log_before"] {
		WriteLog(fmt.Sprintf("%s\n", item))
	}
}

// @func writeLogAfter
// @desc 写入日志后置内容
// @param void
// @return void
func WriteLogAfter() {
	for _, item := range Config["log_after"] {
		WriteLog(fmt.Sprintf("%s\n", item))
	}
}

// @func readProjectLog
// @desc 读取项目日志
// @param project string
// @returns void
func ReadProjectLog(project string) {
	name := path.Base(project)
	// 日志
	gitLogArgs := []string{
		"log",
		"--oneline",
		"--no-merges",
		"--all",
		"--after",
		after + " 00:00:00",
		"--before",
		before + " 23:59:59",
	}

	// 指定author
	for _, author := range Config["author"] {
		gitLogArgs = append(gitLogArgs, "--author", author)
	}

	cmd := exec.Command("git", gitLogArgs...)
	cmd.Dir = project
	if Debug {
		DebugLog(strings.Join(append([]string{path.Join(project) + ": git"}, gitLogArgs...), " "))
	}

	temp, _ := cmd.Output()
	logs := strings.Split(string(temp), "\n")
	temp = nil

	// 倒序
	for i, j := 0, len(logs)-1; i < j; i, j = i+1, j-1 {
		logs[i], logs[j] = logs[j], logs[i]
	}

	// 防止重复
	history := make(map[string]uint, len(logs))

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

			// 重复项过滤
			if history[item] > 0 {
				history[item]++
				continue
			}
			history[item] = 1

			// 判断合并类型
			if Merge {
				MergeLog(name, item)
			} else {
				WriteLog(fmt.Sprintf("%3d、[%v] %v\n", logIdenxNO, name, item))
				logIdenxNO++
			}
		}
	}
}

// @func writeLog
// @desc 写入日志
// @param name string 项目名称
// @param log string 日志文本
// @returns void
func WriteLog(log string) {
	for _, file := range logFiles {
		file.WriteString(log)
	}
}

// @func afterCommand
// @desc 后置操作
// @param void
// @return void
func AfterCommand() {
	for _, cmd := range Config["after_cmd"] {
		for _, path := range logPaths {
			exec.Command(cmd, path).Run()
		}
	}
}

// @func mergeLog
// @desc 按合并规则格式化日志
// @param name string 项目名称
// @param logs string 日志文本
// @return void
func MergeLog(name, logs string) {
	// 遍历规则
	for key, callbackRules := range rulesMap {
		count := 0
		text := ""
		history := []string{}
		for _, callback := range callbackRules {
			if status, _tmp := callback(name, logs); status {
				count++
				// 重复项过滤
				if InSlice(_tmp, history) {
					return
				}
				text = _tmp
				history = append(history, text)
			}
		}
		// 符合规则，合并日志
		if count == len(callbackRules) && text != "" {
			mergeLogs[key] = append(mergeLogs[key], text)
			return
		}
	}
	// 不符合规则，直接写入
	mergeLogs[name] = append(mergeLogs[name], logs)
}

// @func writeMergeLogs
// @desc 写入合并后的日志
// @param void
// @return void
func WriteMergeLogs() {
	for name, logs := range mergeLogs {
		isRule := false
		// 符合写入规则均按规则写入
		for key, callback := range writeMap {
			if name == key {
				callback(name, logs)
				isRule = true
				break
			}
		}

		// 已按规则写入，跳过
		if isRule {
			continue
		}

		// 其它类型日志
		for _, log := range logs {
			log = fmt.Sprintf("%3d、[%v] %v\n", logIdenxNO, name, log)
			WriteLog(log)
			logIdenxNO++
		}
	}
}
