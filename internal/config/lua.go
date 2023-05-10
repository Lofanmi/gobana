package config

import (
	lua "github.com/yuin/gopher-lua"
)

var exports = map[string]lua.LGFunction{
	"gobana_nginx_decode": luaNginxDecode,
}

// GetLuaState
// @autowire(set=config)
func GetLuaState() (L *lua.LState, fn func()) {
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
