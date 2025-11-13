package svc_goja

import (
	"context"

	"github.com/Lofanmi/gobana/service"
	"github.com/dop251/goja"
)

var (
	_ service.GoJa = &Service{}
)

// Service
// @autowire(service.GoJa,set=service)
type Service struct {
	QQWry service.QQWry
}

func (s *Service) GetRuntime() (vm *goja.Runtime) {
	vm = goja.New()
	s.registerServiceToVM(vm)
	return
}

func (s *Service) registerServiceToVM(vm *goja.Runtime) {
	_ = vm.Set("gobana_nginx_decode", s.gobanaNginxDecode(vm))
	_ = vm.Set("gobana_ip_location", s.gobanaIPLocation(vm))
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
