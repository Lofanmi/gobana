package gotil

import (
	"time"

	"github.com/spf13/cast"
	lua "github.com/yuin/gopher-lua"
)

func ParseTime(s string) (timestamp int64) {
	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		return
	}
	timestamp = t.UnixMilli()
	return
}

func MapToTable(m map[string]interface{}) *lua.LTable {
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
		case map[string]interface{}:
			t := MapToTable(res)
			resultTable.RawSetString(key, t)
		case []interface{}:
			sliceTable := &lua.LTable{}
			for _, s := range res {
				switch res2 := s.(type) {
				case map[string]interface{}:
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
