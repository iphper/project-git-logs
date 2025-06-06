package config

import (
	"encoding/json"
	"fmt"
	"os"
	"project-git-logs/cmd"
	"project-git-logs/dialog"
	"strings"
)

type Config struct {
	Path    string `json:"paths"`    // 仓库路径
	LogPath string `json:"log_path"` // 日志存储路径
	Open    bool   `json:"open"`     // 是否开启
	User    string `json:"user"`     // 用户名
}

func (c *Config) Init() {
	// 初始化配置
	c.Path = ""
	c.LogPath = "./logs"
	c.Open = true

	// 读取当前目录下的配置文件
	c.readConfig()
}

// @func readConfig 读取配置文件
func (c *Config) readConfig() {
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			c.initConfig()
		} else {
			panic("读取配置文件失败：" + err.Error())
		}
	} else {

		err = json.Unmarshal(data, c)
		if err != nil {
			panic("解析配置文件失败：" + err.Error())
		}
	}
}

// @func initConfig 初始化配置文件
func (c *Config) initConfig() {
	file, err := os.Create(filename)
	if err != nil {
		panic("创建配置文件失败：" + err.Error())
	}
	defer file.Close()

	// 选择读取目录
	c.Path = dialog.SelectDir("请选择要读取的git仓库目录")

	// 选择日志存储路径
	c.LogPath = dialog.SelectDir("请选择保存日志存储目录")

	// 输入用户名
	c.User = strings.ReplaceAll(cmd.Single().Run("git", "config", "--global", "user.name"), "\n", "")

	// 输入读取完成后是否打开日志
	isOpenLog := "y"
	fmt.Print("是否打开日志？(y/n): ")
	fmt.Scan(&isOpenLog)
	if isOpenLog == "n" || isOpenLog == "N" {
		c.Open = false // 关闭日志
	} else {
		c.Open = true // 打开日志
	}

	jsonData, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		panic("序列化配置文件失败：" + err.Error())
	}

	// 写入到配置文件
	file.Write(jsonData)

}

// ======< getter and setter >======

// @func GetLogPath 获取日志存储路径
func (c *Config) GetLogPath() string {
	return c.LogPath
}

// @func GetPath 获取仓库路径
func (c *Config) GetPath() string {
	return c.Path
}

// @func IsOpen 获取是否开启日志
func (c *Config) IsOpen() bool {
	return c.Open
}
