package utils

import (
	"time"
)

// DateToLocalRange 将 日期字符串(yyyy-MM-dd) 转为 Local 区间字符串
// 输入: "2026-05-26"
// 输出: start(0点对应本地时间), end(UTC 23:59:59对应本地时间)
func DateToLocalRange(startTimeStr, endTimeStr string) (startLocal, endLocal string, err error) {
	loc := getLocalTimeZone()
	// 1. 拼接本地时间: 当天 00:00:00
	localStartStr := startTimeStr + " 00:00:00"
	localStart, err := time.ParseInLocation("2006-01-02 15:04:05", localStartStr, loc)
	if err != nil {
		return "", "", err
	}
	// 2. 拼接本地时间: 当天 23:59:59
	localEndStr := endTimeStr + " 23:59:59"
	localEnd, err := time.ParseInLocation("2006-01-02 15:04:05", localEndStr, loc)
	if err != nil {
		return "", "", err
	}

	// 3. 转为 UTC 并格式化为字符串
	startLocal = localStart.Local().Format("2006-01-02 15:04:05")
	endLocal = localEnd.Local().Format("2006-01-02 15:04:05")
	return startLocal, endLocal, nil
}

func getLocalTimeZone() *time.Location {
	return time.Local
}
