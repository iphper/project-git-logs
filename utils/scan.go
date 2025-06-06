package utils

import "fmt"

// Scan 扫描输入
func Scan(label string, defs ...any) any {
	fmt.Printf("%v", label)
	var val any
	_, err := fmt.Scanf("%v", &val)
	if err != nil || len(val.(string)) == 0 {
		return defs[0] // 如果扫描失败，递归调用
	}
	fmt.Scan()
	return val
}
