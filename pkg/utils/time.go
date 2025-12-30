package utils

import (
	"time"
)

const (
	// TimeLayout 时间格式
	TimeLayout = "2006-01-02 15:04:05"
	// DateLayout 日期格式
	DateLayout = "2006-01-02"
	// TimeLayoutCompact 紧凑时间格式
	TimeLayoutCompact = "20060102150405"
)

// Now 当前时间
func Now() time.Time {
	return time.Now()
}

// NowUnix 当前时间戳（秒）
func NowUnix() int64 {
	return time.Now().Unix()
}

// NowUnixNano 当前时间戳（纳秒）
func NowUnixNano() int64 {
	return time.Now().UnixNano()
}

// NowUnixMilli 当前时间戳（毫秒）
func NowUnixMilli() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// FormatTime 格式化时间
func FormatTime(t time.Time) string {
	return t.Format(TimeLayout)
}

// FormatDate 格式化日期
func FormatDate(t time.Time) string {
	return t.Format(DateLayout)
}

// ParseTime 解析时间字符串
func ParseTime(s string) (time.Time, error) {
	return time.Parse(TimeLayout, s)
}

// ParseDate 解析日期字符串
func ParseDate(s string) (time.Time, error) {
	return time.Parse(DateLayout, s)
}

// UnixToTime 时间戳转时间
func UnixToTime(ts int64) time.Time {
	return time.Unix(ts, 0)
}

// TimeToUnix 时间转时间戳
func TimeToUnix(t time.Time) int64 {
	return t.Unix()
}

// Today 今天开始时间
func Today() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}

// TodayEnd 今天结束时间
func TodayEnd() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location())
}

// WeekStart 本周开始时间（周一）
func WeekStart() time.Time {
	now := time.Now()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	days := weekday - 1
	return time.Date(now.Year(), now.Month(), now.Day()-days, 0, 0, 0, 0, now.Location())
}

// WeekEnd 本周结束时间（周日）
func WeekEnd() time.Time {
	now := time.Now()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	days := 7 - weekday
	return time.Date(now.Year(), now.Month(), now.Day()+days, 23, 59, 59, 999999999, now.Location())
}

// MonthStart 本月开始时间
func MonthStart() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
}

// MonthEnd 本月结束时间
func MonthEnd() time.Time {
	now := time.Now()
	nextMonth := now.Month() + 1
	if nextMonth > 12 {
		nextMonth = 1
	}
	return time.Date(now.Year(), nextMonth, 1, 0, 0, 0, 0, now.Location()).Add(-time.Second)
}

// YearStart 本年开始时间
func YearStart() time.Time {
	now := time.Now()
	return time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
}

// YearEnd 本年结束时间
func YearEnd() time.Time {
	now := time.Now()
	return time.Date(now.Year(), 12, 31, 23, 59, 59, 999999999, now.Location())
}

// AddDays 添加天数
func AddDays(t time.Time, days int) time.Time {
	return t.AddDate(0, 0, days)
}

// AddMonths 添加月数
func AddMonths(t time.Time, months int) time.Time {
	return t.AddDate(0, months, 0)
}

// AddYears 添加年数
func AddYears(t time.Time, years int) time.Time {
	return t.AddDate(years, 0, 0)
}

// DiffDays 计算两个时间相差的天数
func DiffDays(t1, t2 time.Time) int {
	return int(t1.Sub(t2).Hours() / 24)
}

// IsSameDay 判断是否为同一天
func IsSameDay(t1, t2 time.Time) bool {
	return t1.Year() == t2.Year() && t1.Month() == t2.Month() && t1.Day() == t2.Day()
}
