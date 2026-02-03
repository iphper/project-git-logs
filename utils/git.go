package utils

import (
	"os"
	"path/filepath"
	"strings"
)

func GitGlobalUserName() string {
	// git config --global user.name
	return Cmd("git", "config", "--global", "user.name")
}

func GitRepositories(dirname string) []string {
	entries, err := os.ReadDir(dirname)
	if err != nil {
		panic(err)
	}

	var repos []string

	// 先尝试本身是否是git仓库
	gitPath := filepath.Join(dirname, ".git")
	if _, err := os.Stat(gitPath); err == nil {
		return []string{dirname}
	}

	// 再尝试子目录
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// 可选：跳过超大目录
		switch entry.Name() {
		case "node_modules", "vendor", ".cache":
			continue
		}

		gitPath := filepath.Join(dirname, entry.Name(), ".git")
		if _, err := os.Stat(gitPath); err == nil {
			repos = append(repos, filepath.Join(dirname, entry.Name()))
		} else {
			repos = append(repos, GitRepositories(filepath.Join(dirname, entry.Name()))...)
		}
	}
	return repos
}

func GitFormat(line string) string {
	prefixs := []string{
		"feat:", "fix:", "perf:", "style:", "docs:", "test:", "refactor:", "build:", "ci:", "chore:", "revert:",
		"wip:", "workflow:", "types:", "release:",
	}
	line = strings.TrimSpace(line)

	for _, prefix := range prefixs {
		if strings.HasPrefix(line, prefix) {
			line = strings.TrimPrefix(line, prefix)
			line = strings.TrimSpace(line)
			break
		}
	}
	return line
}
