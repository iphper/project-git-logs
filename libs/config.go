package libs

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// 常量
const (
	// 获取配置文件【保存一些配置到此文件】
	// 如果存在此文件就获取此文件中的参数
	// 如果不存在初始化时需要使用者输入
	configFileName = "weeker.conf"
)

var (
	// 是否开启debug
	Debug bool = false
	// 是否合并日志
	Merge = true
	// 配置数据
	Config map[string][]string = make(map[string][]string, 1)
)

// @func InitArgs
// @desc 初始化参数
// @param void
// @return void
func InitArgs() {
	// 遍历命令参数
	for _, arg := range os.Args {
		// 默认开启debug
		if arg == "-debug" || arg == "-d" {
			Debug = true
			break
		}
		// 是否合并提交
		if arg == "-no-merges" || arg == "-nms" {
			Merge = false
			break
		}
	}
}

// @func initConfig
// @desc 初始化配置
// @param void
// @return void
func InitConfig() {
	// 输入日期范围
	InputTimes()

	// 配置参数获取
	if !PathExists(configFileName) {
		InitConfigFile()
	} else {
		ReadConfig()
	}

	// 配置占位符替换
	ConfigPlaceholder()
}

// @func initConfigFileInitConfigFile
// @desc 初始化配置
// @param void
// @return void
func InitConfigFile() {
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

		Config[field] = values[:i-1]
	}

	// 询问是否写入配置文件
	for {
		fmt.Print("是否保存到配置文件？[yes|no]:")
		fmt.Scanln(&input)

		input = strings.ToLower(input)
		if input == "yes" || input == "y" {
			// 保存到配置文件
			WriteConfig()
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
func ReadConfig() {
	var err error
	// 存在配置文件就读取配置文件内容
	bufs, e := os.ReadFile(configFileName)
	if e != nil {
		HandleError(err)
	}

	// 读取配置
	err = json.Unmarshal(bufs, &Config)
	if err != nil {
		HandleError(err)
	}
}

// @func writeConfig
// @desc 写入配置文件
// @param void
// @return void
func WriteConfig() {
	file, err := os.Create(configFileName)
	if err != nil {
		HandleError(err)
	}

	defer file.Close()

	bufs, err := json.Marshal(Config)
	if err != nil {
		HandleError(err)
	}

	file.Write(bufs)
}

// @func getProjectList
// @desc 获取项目列表
// @param void
// @return []string
func GetProjectList() map[string]uint {

	// 获取所有的项目目录列表
	projects := map[string]uint{}
	for _, item := range Config["project_root"] {
		for _, project := range LoadProjects(item) {
			projects[project] = 1
		}
	}

	return projects
}

// @func configPlaceholder
// @desc 配置占位符转化
// @param void
// @return void
func ConfigPlaceholder() {

	// 配置占位符转化
	for field, values := range Config {
		for idx, value := range values {
			switch true {
			// 时间日期点位符转化
			case strings.Contains(value, "$DATE"):
				Config[field][idx] = strings.ReplaceAll(value, "$DATE", DateFormat())
			case strings.Contains(value, "$TIME"):
				Config[field][idx] = strings.ReplaceAll(value, "$DATE", TimeFormat())
			case strings.Contains(value, "$DATETIME"):
				Config[field][idx] = strings.ReplaceAll(value, "$DATE", DateTimeFormat())
			}
		}
	}
}

// @func inputTimes
// @desc 输入开始和结束时间
// @param void
// @return void
func InputTimes() {

	var input string
	// 获取获取开始日期
	for {
		fmt.Printf("请输入读取提交日志的开始时间[%v]:", after)
		fmt.Scanln(&input)
		if input != "" {
			if IsValidDate(input) {
				after = input[:10]
				break
			}
			fmt.Println("输入的日期格式不正确")
		} else {
			break
		}
	}

	// 重新初始化
	input = ""
	// 获取获取结束日期
	for {
		fmt.Printf("请输入读取提交日志的结束时间[%v]:", before)
		fmt.Scanln(&input)
		if input != "" {
			if IsValidDate(input) {
				before = input[:10]
				break
			}
			fmt.Println("输入的日期格式不正确")
		} else {
			break
		}
	}
}
