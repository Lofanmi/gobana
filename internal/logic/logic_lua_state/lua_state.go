package logic_lua_state

import (
	"encoding/json"

	"github.com/Lofanmi/gobana/internal/gotil"
	"github.com/Lofanmi/gobana/internal/logic"
	lua "github.com/yuin/gopher-lua"
)

var (
	_       logic.LuaState = &LuaState{}
	exports                = map[string]lua.LGFunction{
		"gobana_nginx_decode": luaNginxDecode,
		"gobana_json_decode":  luaJsonDecode,
	}
)

// LuaState
// @autowire(logic.LuaState,set=logics)
type LuaState struct {
}

func (s LuaState) GetLuaState() (L *lua.LState, fn func()) {
	L = lua.NewState()
	fn = func() { L.Close() }
	for name, function := range exports {
		L.SetGlobal(name, L.NewFunction(function))
	}
	return
}

func luaNginxDecode(L *lua.LState) int {
	s := L.ToString(1)
	L.Push(lua.LString(nginxDecode(s)))
	return 1
}

func luaJsonDecode(L *lua.LState) int {
	s := L.ToString(1)
	var m map[string]interface{}
	_ = json.Unmarshal([]byte(s), &m)
	L.Push(gotil.MapToTable(m))
	return 1
}
