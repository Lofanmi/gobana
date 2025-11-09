package svc_lua_state

import (
	"context"

	"github.com/Lofanmi/gobana/service"
	lua "github.com/yuin/gopher-lua"
)

var (
	_       service.LuaState = &LuaState{}
	exports                  = map[string]func(state *LuaState) lua.LGFunction{
		"gobana_nginx_decode": luaNginxDecode,
		"gobana_ip_location":  luaIPLocation,
	}
)

// LuaState
// @autowire(service.LuaState,set=service)
type LuaState struct {
	QQWry service.QQWry
}

func (s *LuaState) GetLuaState() (L *lua.LState, fn func()) {
	L = lua.NewState()
	fn = func() { L.Close() }
	for name, builder := range exports {
		L.SetGlobal(name, L.NewFunction(builder(s)))
	}
	return
}

func luaNginxDecode(state *LuaState) lua.LGFunction {
	_ = state
	return func(L *lua.LState) int {
		s := L.ToString(1)
		L.Push(lua.LString(nginxDecode(s)))
		return 1
	}
}

func luaIPLocation(state *LuaState) lua.LGFunction {
	return func(L *lua.LState) int {
		var res string
		ip := L.ToString(1)
		location, err := state.QQWry.Find(context.Background(), ip)
		if err == nil {
			res = location.String()
		}
		L.Push(lua.LString(res))
		return 1
	}
}
