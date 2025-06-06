package git

import (
	"os"
	"path/filepath"
	"sync"
)

var (
	log  *Log
	once sync.Once
)

// @func init 初始化
func init() {
	once.Do(func() {
		log = &Log{}
		log.Init()
	})
}

// @func ReadGitPaths 读取指定路径下的所有git仓库路径
func ReadGitPaths(p string) []string {

	// 读取指定路径下的所有git仓库路径
	// 这里可以使用os和filepath包来遍历目录
	var gitPaths []string

	// 当前路径是否是一个git仓库
	if IsGitPath(p) {
		return []string{p}
	}

	// 遍历目录
	dirs, err := os.ReadDir(p)
	if err != nil {
		panic(err)
	}
	for _, dir := range dirs {
		if dir.IsDir() {
			subPath := filepath.Join(p, dir.Name())
			if IsGitPath(subPath) {
				gitPaths = append(gitPaths, subPath)
			} else {
				// 递归查找子目录
				gitPaths = append(gitPaths, ReadGitPaths(subPath)...)
			}
		}
	}

	return gitPaths
}

// @func IsGitPath 判断给定路径是否是一个git仓库
func IsGitPath(path string) bool {
	// 判断给定路径是否是一个git仓库
	_, err := os.Stat(filepath.Join(path, ".git"))
	if err == nil {
		return true
	}
	// 这里可以检查路径下是否存在.git目录
	return false // 这里需要实现具体的逻辑
}

// @func ReadGitLogs 读取指定路径下的git仓库提交日志
func ReadGitLogs(path string) error {
	if !IsGitPath(path) {
		return nil
	}

	return log.Read(path)
}

// @func WriteGitLogs 写入git仓库提交日志
func WriteGitLogs(f string) {
	log.Write(f)
}
