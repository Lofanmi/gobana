package gotil

import (
	"strconv"
	"strings"
	"time"

	"github.com/hako/durafmt"
)

func FormatDuration(duration time.Duration) string {
	ret := strings.Replace(durafmt.Parse(duration).LimitFirstN(3).InternationalString(), " ", "", -1)
	if strings.HasSuffix(ret, "1") {
		ret += "s"
	}
	return ret
}

func Ymd() (i int) {
	s := time.Now().Format("20060102")
	i, _ = strconv.Atoi(s)
	return
}

func DateMs(ts int64) (s string) {
	s = time.UnixMilli(ts).Format("2006-01-02 15:04:05")
	return
}

func DateSec(ts int64) (s string) {
	s = time.Unix(ts, 0).Format("2006-01-02 15:04:05")
	return
}

func IfElse[T any](condition bool, a, b T) T {
	if condition {
		return a
	}
	return b
}
