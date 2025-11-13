package gotil

import (
	"time"
)

func FormatMilliSecond(ts int64) (s string) {
	s = time.UnixMilli(ts).Format("2006-01-02 15:04:05")
	return
}

func ParseTime(s string) (timestamp int64) {
	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		return
	}
	timestamp = t.UnixMilli()
	return
}

func DefaultInterval(timeA, timeB, maxChartPoints int64) (interval int) {
	interval = (int)((timeB - timeA) / 1000 / maxChartPoints)
	if interval <= 1 {
		interval = 1
	} else if interval <= 5 {
		interval = 5
	} else if interval <= 10 {
		interval = 10
	} else if interval <= 30 {
		interval = 30
	} else if interval <= 60 {
		interval = 60
	} else if interval <= 300 {
		interval = 300
	} else if interval <= 900 {
		interval = 900
	} else if interval <= 1800 {
		interval = 1800
	} else if interval <= 3600 {
		interval = 3600
	} else if interval <= 3600*3 {
		interval = 3600 * 3
	} else if interval <= 3600*9 {
		interval = 3600 * 9
	} else if interval <= 3600*12 {
		interval = 3600 * 12
	} else {
		interval = 3600 * 24
	}
	return
}
