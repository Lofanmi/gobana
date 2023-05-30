package app

import (
	"github.com/Lofanmi/gobana/internal/config"
	"github.com/Lofanmi/gobana/service"
	"github.com/gin-gonic/gin"
)

const (
	ServerPrefix = "/api/gobana/v1"
)

// Application
// @autowire(set=app)
type Application struct {
	ConfigApplication config.Application
	Config            service.Config
	Logger            service.Logger
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (s *Application) Run() {
	if !s.ConfigApplication.Production {
		GenStructGraph(s, "design/structure.svg")
	}
	engine := gin.Default()
	s.registerStaticRouter(engine)
	s.registerApiRouter(engine.Group(ServerPrefix))
	if err := engine.Run(s.ConfigApplication.ListenAddr); err != nil {
		panic(err)
	}
}
