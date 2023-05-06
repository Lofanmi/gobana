// Code generated by go-autowire. DO NOT EDIT.

//go:build wireinject
// +build wireinject

package inject

import (
	"github.com/google/wire"

	"github.com/Lofanmi/gobana/internal/logic"
	"github.com/Lofanmi/gobana/internal/logic/logic_backend_factory"
	"github.com/Lofanmi/gobana/internal/logic/logic_query_builder"
)

var LogicsSet = wire.NewSet(
	logic_backend_factory.NewBackendFactory,

	wire.Struct(new(logic_query_builder.QueryBuilder), "*"),
	wire.Bind(new(logic.QueryBuilder), new(*logic_query_builder.QueryBuilder)),
)