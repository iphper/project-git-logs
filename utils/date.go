package utils

import "time"

// GetWeekStartDate 返回本周第一天（周一）的日期字符串，格式为"2006-01-02"
func GetWeekStartDate() string {
	now := time.Now()
	weekday := int(now.Weekday())
	if weekday == 0 { // 周日
		weekday = 7
	}
	monday := now.AddDate(0, 0, -weekday+1)
	return monday.Format("2006-01-02")
}

// GetTodayDate 返回今天的日期字符串，格式为"2006-01-02"
func GetTodayDate() string {
	return time.Now().Format("2006-01-02")
}

// IsValidDate 判断日期字符串是否合法，格式应为"2006-01-02"
func IsValidDate(dateStr string) bool {
	_, err := time.Parse("2006-01-02", dateStr)
	return err == nil
}
