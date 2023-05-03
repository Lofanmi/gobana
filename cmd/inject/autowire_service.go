// Code generated by go-autowire. DO NOT EDIT.

//go:build wireinject
// +build wireinject

package inject

import (
	"github.com/google/wire"

	"github.com/Lofanmi/gobana/internal/svc_logger"
	"github.com/Lofanmi/gobana/service"
)

var ServiceSet = wire.NewSet(
	svc_logger.NewService,
	wire.Bind(new(service.Logger), new(*svc_logger.Service)),
)
