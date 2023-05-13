package logic

import (
	lua "github.com/yuin/gopher-lua"
)

type LuaState interface {
	GetLuaState() (L *lua.LState, fn func())
}
