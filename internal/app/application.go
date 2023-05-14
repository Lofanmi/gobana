package app

import (
	"context"
	"fmt"

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
	engine := gin.Default()
	s.registerStaticRouter(engine)
	s.registerApiRouter(engine.Group(ServerPrefix))
	if err := engine.Run(s.ConfigApplication.ListenAddr); err != nil {
		panic(err)
	}
}

func (s *Application) RunTest() {
	ctx := context.Background()
	var q service.QueryByHuman
	q.Must = append(q.Must, "uid")
	resp, err := s.Logger.Search(ctx, service.SearchRequest{
		PageNo:   1,
		PageSize: 10,
		Backend:  "全球平台",
		Storage:  "all",
		QueryBy:  "query_by_human",
		Query:    &q,
	})
	fmt.Printf("%+v\n", resp)
	fmt.Printf("%+v\n", err)
}
