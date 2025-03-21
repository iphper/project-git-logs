package src

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Cache struct{}

var (
	cache *Cache
)

func GetCache() *Cache {
	cache = new(Cache)
	cache.Initialite()
	return cache
}

func (c *Cache) Initialite() {

}

func (c *Cache) Read(filename string) map[string][]string {
	// 读取数据
	data := map[string][]string{}
	var err error
	// 存在配置文件就读取配置文件内容
	bufs, e := os.ReadFile(filename)
	if e != nil {
		panic(e)
	}

	// 读取配置
	err = json.Unmarshal(bufs, &data)
	if err != nil {
		panic(err)
	}
	return data
}

func (c *Cache) Write(configFileName string, config map[string][]string) {
	file, err := os.Create(configFileName)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	bufs, err := json.Marshal(config)
	if err != nil {
		panic(err)
	}

	file.Write(bufs)
}

func (c *Cache) Create(file string) bool {
	fmt.Println("-----------------------------")
	fmt.Println("未找到配置信息，请填写以下配置项")

	keys := map[string]string{
		"root":  "项目路径",
		"after": "完成后执行命令操作",
	}

	var input string
	var values []string
	var i int

	config := map[string][]string{}

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
			cache.Write(file, config)
			break
		}
		if input == "no" || input == "n" {
			// 不保存
			break
		}
	}
	fmt.Println("-----------------------------")

	return true
}
