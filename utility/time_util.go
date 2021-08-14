package utility

import (
	"strconv"
	"time"
)

func ParseStrToTimestamp(str string) (int64, error) {
	intTime, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}
	return int64(intTime), nil
}

func FormatTime(timestamp int64) string {
	return time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")
}

func GetNowTimeStr() string {
	return FormatTime(time.Now().Local().Unix())
}

func GetCurrentMillisecond() int64 {
	return time.Now().UnixNano() / (1000 * 1000)
}
