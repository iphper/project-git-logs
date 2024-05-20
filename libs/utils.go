package libs

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"
)

// @func inSlice
// @desc 判断是否在切片中
// @param item string
// @param list []string
// @return bool
func InSlice(item string, list []string) bool {
	for _, v := range list {
		if item == v {
			return true
		}
	}
	return false
}

// @func isValiDate
// @desc 判断是否为合法的日期格式
// @param dateStr string 日期字符串
// @return bool
func IsValidDate(dateStr string) bool {
	_, err := time.Parse("2006-01-02", dateStr)
	return err == nil
}

// @func debugLog
// @desc 调试日志
// @param msg []string
// @return void
func DebugLog(msg ...string) {
	if !Debug {
		return
	}

	nowDir, _ := os.Getwd()
	file := path.Join(nowDir, "debug.log")
	fpointer, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		HandleError(err)
	}

	defer fpointer.Close()

	for _, m := range msg {
		m = fmt.Sprintf("[%s] %s\n", time.Now().Format("2006-01-02 15:04:05"), m)
		fpointer.WriteString(m)
	}
}

// @func handleError
// @desc 处理错误
// @param err error
// @return void
func HandleError(err error) {
	fmt.Println(err.Error())
	os.Exit(1)
}

// @func isProject
// @desc 是否为项目目录
// @param name string
// @return bool
func IsProject(name string) bool {
	dirs, _ := os.ReadDir(name)
	str := ".git,README.md,.idea,"
	for _, item := range dirs {
		if strings.Contains(str, item.Name()+",") {
			return true
		}
	}
	return false
}

// @func loadProjects
// @desc 加载项目目录
// @param projectPath string
// @return []string
func LoadProjects(projectPath string) (list []string) {
	if IsProject(projectPath) {
		list = append(list, projectPath)
	} else {
		dirs, _ := os.ReadDir(projectPath)
		var temp string

		for _, project := range dirs {
			temp = path.Join(projectPath, project.Name())
			if IsProject(temp) {
				list = append(list, temp)
				continue
			}
			if project.IsDir() {
				list = append(list, LoadProjects(temp)...)
			}
		}
	}
	return
}

// @func pathExists
// @desc 路径是否存在
// @param path string
// @return bool
func PathExists(path string) bool {
	if len(path) < 1 {
		return false
	}

	_, e := os.Stat(path)
	if e == nil {
		return true
	}
	return os.IsExist(e)
}

// @func DateFormat
// @desc 格式化日期
// @param timestamp ...time.Time
// @return string
func DateFormat(timestamp ...time.Time) string {
	if len(timestamp) == 0 {
		return time.Now().Format("2006-01-02")
	}
	dateArr := []string{}
	for _, item := range timestamp {
		dateArr = append(dateArr, item.Format("2006-01-02"))
	}
	return strings.Join(dateArr, ",")
}

// @func TimeFormat
// @desc 格式化时间
// @param timestamp ...time.Time
// @return string
func TimeFormat(timestamp ...time.Time) string {
	if len(timestamp) == 0 {
		return time.Now().Format("15:04:05")
	}
	dateArr := []string{}
	for _, item := range timestamp {
		dateArr = append(dateArr, item.Format("15:04:05"))
	}
	return strings.Join(dateArr, ",")
}

// @func DateTimeFormat
// @desc 格式化日期时间
// @param timestamp ...time.Time
// @return string
func DateTimeFormat(timestamp ...time.Time) string {
	if len(timestamp) == 0 {
		return time.Now().Format("2006-01-02 15:04:05")
	}
	dateArr := []string{}
	for _, item := range timestamp {
		dateArr = append(dateArr, item.Format("2006-01-02 15:04:05"))
	}
	return strings.Join(dateArr, ",")
}
