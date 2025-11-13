package gotil

import (
	"time"

	"github.com/spf13/cast"
	lua "github.com/yuin/gopher-lua"
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

func MapToTable(m map[string]any) *lua.LTable {
	resultTable := &lua.LTable{}
	for key, element := range m {
		switch res := element.(type) {
		case float64, float32, int64, int32, int16, int8, uint, uint64, uint32, uint16, uint8:
			resultTable.RawSetString(key, lua.LNumber(cast.ToFloat64(res)))
		case string:
			resultTable.RawSetString(key, lua.LString(res))
		case bool:
			resultTable.RawSetString(key, lua.LBool(res))
		case []byte:
			resultTable.RawSetString(key, lua.LString(res))
		case map[string]any:
			t := MapToTable(res)
			resultTable.RawSetString(key, t)
		case []any:
			sliceTable := &lua.LTable{}
			for _, s := range res {
				switch res2 := s.(type) {
				case map[string]any:
					t := MapToTable(res2)
					sliceTable.Append(t)
				case float64, float32, int64, int32, int16, int8, uint, uint64, uint32, uint16, uint8:
					resultTable.RawSetString(key, lua.LNumber(cast.ToFloat64(res2)))
				case string:
					resultTable.RawSetString(key, lua.LString(res2))
				case bool:
					resultTable.RawSetString(key, lua.LBool(res2))
				case []byte:
					resultTable.RawSetString(key, lua.LString(res2))
				}
			}
			resultTable.RawSetString(key, sliceTable)
		}
	}
	return resultTable
}
