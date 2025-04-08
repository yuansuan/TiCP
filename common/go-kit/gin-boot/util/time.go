package util

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/yuansuan/ticp/common/go-kit/logging"
)

// FirstOfMonth get the start timestamp of current month
func FirstOfMonth() time.Time {
	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	return time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, time.Local)
}

// ParseTimeDuration parses strings like 1-00:00:00 or 01:23:55
func ParseTimeDuration(t string) (time.Duration, error) {
	parts := strings.Split(t, "-")
	var days int
	var tm string
	if len(parts) == 1 {
		tm = parts[0]
	} else {
		tm = parts[1]
		d, err := strconv.Atoi(parts[0])
		if err != nil {
			return 0, fmt.Errorf("parseTimeDuration[%v] format error: get day %v", t, err)
		}
		days = d
	}

	ts := strings.Split(tm, ":")
	if len(ts) != 3 {
		return 0, fmt.Errorf("parseTimeDuration[%v] format error: hh:MM:ss", t)
	}
	h, err := strconv.Atoi(ts[0])
	if err != nil {
		return 0, fmt.Errorf("parseTimeDuration[%v] format error: get h %v", t, err)
	}
	M, err := strconv.Atoi(ts[1])
	if err != nil {
		return 0, fmt.Errorf("parseTimeDuration[%v] format error: get M %v", t, err)
	}
	s, err := strconv.Atoi(ts[2])
	if err != nil {
		return 0, fmt.Errorf("parseTimeDuration[%v] format error: get s %v", t, err)
	}
	return time.Duration(days)*24*time.Hour + time.Duration(h)*time.Hour + time.Duration(M)*time.Minute + time.Duration(s)*time.Second, nil
}

// TimeTrack measures time consuming
func TimeTrack(start time.Time, ctx string) {
	elapsed := time.Since(start)
	logging.Default().Debugf("%v consumes %v", ctx, elapsed.String())
}
