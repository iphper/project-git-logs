package dialog

import (
	"github.com/gen2brain/dlgs"
)

// 选择目录
func SelectDir(title string) string {
	dir, ok, err := dlgs.File(title, "", true) // `true` 表示选择目录
	if err != nil {
		panic(title + "错误：" + err.Error())
	}

	// 用户取消
	if !ok {
		return ""
	}

	return dir
}
