package service

import (
	"github.com/dop251/goja"
)

type GoJa interface {
	GetRuntime() (vm *goja.Runtime)
}
