package app

import (
	"context"
	"fmt"

	"github.com/Lofanmi/gobana/service"
)

// Application
// @autowire(set=app)
type Application struct {
	Logger service.Logger
}

func (s *Application) Run() {
	ctx := context.Background()
	var q service.QueryByHuman
	q.Must = append(q.Must, "lumen")
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
