package svc_goja

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cast"
)

func formatTime(s string) (res string) {
	s = strings.TrimSpace(s)
	if regexp.MustCompile(`^\d+$`).MatchString(s) {
		i, _ := strconv.Atoi(s)
		if len(s) == 10 {
			i *= 1000
		}
		return time.UnixMilli(int64(i)).Format(time.RFC3339Nano)
	}
	if strings.HasSuffix(s, "Z") || strings.Contains(s, "+") {
		t, err := time.ParseInLocation(time.RFC3339, s, time.Local)
		if err == nil {
			return t.Format(time.RFC3339Nano)
		}
	}
	t, err := time.Parse("2006-01-02 15:04:05.000000Z07:00", s)
	if err == nil {
		return t.Format(time.RFC3339Nano)
	}
	t, err = time.Parse("2006-01-02 15:04:05.000Z07:00", s)
	if err == nil {
		return t.Format(time.RFC3339Nano)
	}
	t, err = cast.ToTimeInDefaultLocationE(s, time.Local)
	if err == nil {
		return t.Format(time.RFC3339Nano)
	}
	return s
}

func formatDuration(s string) (res string) {
	duration, err := strconv.ParseFloat(s, 64)
	if err != nil {
		res = s
		return
	}
	t := time.Duration(duration * float64(time.Second))
	res = t.String()
	return
}
