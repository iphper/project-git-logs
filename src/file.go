package src

import (
	"os"
	"path"
	"strings"
)

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
