package timeutil

import (
	"github.com/yuansuan/ticp/PSP/psp/internal/common"
	"github.com/yuansuan/ticp/PSP/psp/test/consts"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"strconv"
	"time"
)

func GetUnixMilli() int64 {
	return time.Now().UnixMilli()
}

func GetUnixMilliStr() string {
	return strconv.FormatInt(GetUnixMilli(), 10)
}

func GetRangeTime(internal int, dateType string) (time.Time, time.Time) {
	currentTime := time.Now()
	agoTime := time.Now()
	switch dateType {
	case consts.YearType:
		agoTime = currentTime.AddDate(-internal, 0, 0)
	case consts.MonthType:
		agoTime = currentTime.AddDate(0, -internal, 0)
	case consts.DayType:
		agoTime = currentTime.AddDate(0, 0, -internal)
	default:
		logging.Default().Warn("date type not match")
	}
	return agoTime, currentTime
}

func GetRangeTimeStr(internal int, dateType string) (string, string) {
	agoTime, currentTime := GetRangeTime(internal, dateType)
	return agoTime.Format(common.DatetimeFormat), currentTime.Format(common.DatetimeFormat)
}
