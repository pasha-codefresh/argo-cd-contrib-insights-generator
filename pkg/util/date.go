package util

import (
	"strconv"
	"time"
)

func GetRangeForLastWeekAsMilli() (string, string) {
	from := strconv.FormatInt(time.Now().AddDate(0, 0, -7).UnixMilli(), 10)
	to := strconv.FormatInt(time.Now().UnixMilli(), 10)
	return from, to
}

func GetRangeForLastWeek() (string, string) {
	from := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	to := time.Now().Format("2006-01-02")
	return from, to
}
