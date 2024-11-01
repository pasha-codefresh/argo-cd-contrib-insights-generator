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
