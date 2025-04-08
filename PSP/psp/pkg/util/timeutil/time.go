package timeutil

import (
	"math"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"

	"github.com/yuansuan/ticp/PSP/psp/internal/common"
)

var (
	DefaultLocation, _ = time.LoadLocation("Asia/Shanghai")
	DefaultDateTime    = time.Date(1970, 1, 1, 8, 0, 0, 0, DefaultLocation)
)

func DefaultFormatTime(t time.Time) string {
	return FormatTime(t, time.RFC3339)
}

func FormatTime(t time.Time, format string) string {
	if t.IsZero() {
		return ""
	}

	locationTime := t.In(DefaultLocation).Format(format)
	if locationTime == DefaultDateTime.Format(format) {
		return ""
	}

	return locationTime
}

func ParseJsonTime(t string) (time.Time, error) {
	return time.ParseInLocation(time.RFC3339, t, DefaultLocation)
}

func ParseTimeWithFormat(timeStr, format string) (time.Time, error) {
	return time.ParseInLocation(format, timeStr, DefaultLocation)
}

func ParseTimeStrWithFormat(timeStr, originFormat, newFormat string) (string, error) {
	newTime, err := ParseTimeWithFormat(timeStr, originFormat)
	if err != nil {
		return "", err
	}
	return newTime.Format(newFormat), nil
}

func ParseUnixTime(t int64) time.Time {
	return time.Unix(t, 0)
}

func ParseTime(t time.Time) int64 {
	return t.Unix()
}

func GetTimeIntervalDay(t1, t2 time.Time) int {
	return int(math.Abs(t1.Sub(t2).Seconds()) / common.OneDayToSecond)
}

func ParseProtoTimeToTime(t *timestamp.Timestamp) time.Time {
	return time.Unix(t.Seconds, int64(t.Nanos))
}

func GetTomorrow() (start, end int64) {
	// 获取当前时间
	now := time.Now()

	// 获取明天的日期
	tomorrow := now.AddDate(0, 0, 1)

	// 构建明天的起始时间和结束时间
	startOfDay := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, tomorrow.Location())
	endOfDay := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 23, 59, 59, 999999999, tomorrow.Location())

	return startOfDay.Unix(), endOfDay.Unix()
}
