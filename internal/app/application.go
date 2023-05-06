package app

import (
	"github.com/Lofanmi/gobana/service"
)

// Application
// @autowire(set=app)
type Application struct {
	Logger service.Logger
}

func (s *Application) Run() {
	//
}
