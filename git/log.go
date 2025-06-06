package git

import (
	"fmt"
	"os"
	"path"
	"project-git-logs/cmd"
	"project-git-logs/utils"
	"strconv"
	"strings"
)

type Log struct {
	content map[string]string
}

// @func Init 初始化日志
func (log *Log) Init() {
	// 初始化日志内容
	log.content = make(map[string]string, 512)
}

// @func Read 读取git提交日志
func (log *Log) Read(p string) error {
	// 读取指定路径下的git提交日志
	content := cmd.Single().PathRun(
		p,
		"git",
		"log",
		"--pretty=format:%s",
		"--no-merges",
		"--author", cmd.Single().GetOption("author").(string),
		"--after", cmd.Single().GetOption("before").(string),
		"--before", cmd.Single().GetOption("after").(string),
	)

	if len(content) == 0 {
		return nil
	}

	// ;替换分号为换行符
	content = strings.ReplaceAll(content, ";", "\n")

	for _, item := range utils.Slice_Duplicates(strings.Split(content, "\n")) {
		// 格式化
		item = strings.TrimSpace(item)
		if len(item) == 0 {
			continue
		}
		item = log.Format(item)

		// 如果内容已经存在，则不重复添加
		if _, ok := log.content[item]; !ok {
			log.content[item] = item
		}
	}

	return nil // 返回nil表示成功
}

// @func Write 写入git提交日志到指定文件
func (log *Log) Write(f string) {
	// 判断目录是否存在
	dir := path.Dir(f)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, 0755)
	}

	fp, err := os.OpenFile(f, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}

	defer fp.Close()

	// 这里需要实现具体的写入git提交日志的逻辑
	// 例如将读取到的日志写入到指定的文件中

	fmat := "%" + strconv.Itoa(len(strconv.Itoa(len(log.content)))) + "d: %v\n"

	i := 1
	for _, line := range log.content {
		line = fmt.Sprintf(fmat, i, line)
		fp.WriteString(line)
		i++
	}
}

// @func Format 格式化提交信息
func (log *Log) Format(line string) string {
	prefixs := []string{
		"feat:", "fix:", "perf:", "style:", "docs:", "test:", "refactor:", "build:", "ci:", "chore:", "revert:",
		"wip:", "workflow:", "types:", "release:",
	}

	for _, prefix := range prefixs {
		if strings.HasPrefix(line, prefix) {
			line = strings.TrimPrefix(line, prefix)
			line = strings.TrimSpace(line)
			break
		}
	}
	// 这里可以实现具体的格式化逻辑
	// 例如将提交信息进行特定的格式化处理
	return line // 返回格式化后的字符串
}
