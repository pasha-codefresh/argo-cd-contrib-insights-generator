package util

import (
	"time"
)

func GetRangeForLastWeek() (string, string) {
	from := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	to := time.Now().Format("2006-01-02")
	return from, to
}
