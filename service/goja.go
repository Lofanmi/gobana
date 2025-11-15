package service

import (
	"github.com/dop251/goja"
)

type GoJa interface {
	GetCallable(parserName, functionName string) (vm *goja.Runtime, callable goja.Callable, err error)
}
