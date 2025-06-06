package config

import "sync"

var (
	config   *Config // 全局配置变量
	once     sync.Once
	filename = "config.json" // 配置文件名
)

// @func Single 获取单例配置
func Single() *Config {
	once.Do(func() {
		config = &Config{}
		config.Init()
	})
	return config
}
