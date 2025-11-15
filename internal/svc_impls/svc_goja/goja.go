package svc_goja

import (
	"context"
	"errors"

	"github.com/Lofanmi/gobana/internal/config"
	"github.com/Lofanmi/gobana/service"
	"github.com/dop251/goja"
)

var _ service.GoJa = &Service{}

// Service
// @autowire(service.GoJa,set=service)
type Service struct {
	Parsers config.Parsers
	QQWry   service.QQWry
}

func (s *Service) registerServiceToVM(vm *goja.Runtime) {
	_ = vm.Set("gobanaNginxDecode", s.gobanaNginxDecode(vm))
	_ = vm.Set("gobanaIPLocation", s.gobanaIPLocation(vm))
	_ = vm.Set("gobanaFormatTime", s.gobanaFormatTime(vm))
	_ = vm.Set("gobanaFormatDuration", s.gobanaDuration(vm))
}

func (s *Service) GetCallable(parserName, functionName string) (vm *goja.Runtime, callable goja.Callable, err error) {
	parser := s.Parsers[parserName]
	if parser.Program == nil {
		err = errors.New("parser program not found")
		return
	}
	vm = goja.New()
	s.registerServiceToVM(vm)
	_, err = vm.RunProgram(parser.Program)
	if err != nil {
		return
	}
	fn := vm.Get(functionName)
	if fn == nil || goja.IsUndefined(fn) || goja.IsNull(fn) {
		err = errors.New("function not found in JavaScript program")
		return
	}
	callable, ok := goja.AssertFunction(fn)
	if !ok {
		err = errors.New("value is not a function")
		return
	}
	return vm, callable, nil
}

func (s *Service) gobanaNginxDecode(vm *goja.Runtime) func(goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		input := call.Argument(0).String()
		return vm.ToValue(nginxDecode(input))
	}
}

func (s *Service) gobanaIPLocation(vm *goja.Runtime) func(goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		ip := call.Argument(0).String()
		location, err := s.QQWry.Find(context.Background(), ip)
		if err != nil {
			return vm.ToValue("")
		}
		return vm.ToValue(location.String())
	}
}

func (s *Service) gobanaFormatTime(vm *goja.Runtime) func(goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		input := call.Argument(0).String()
		return vm.ToValue(formatTime(input))
	}
}

func (s *Service) gobanaDuration(vm *goja.Runtime) func(goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		input := call.Argument(0).String()
		return vm.ToValue(formatDuration(input))
	}
}
