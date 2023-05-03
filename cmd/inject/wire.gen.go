//go:build wireinject

package inject

import (
	"github.com/google/wire"

	"github.com/Lofanmi/gobana/internal/app"
)

func NewApplication() (*app.Application, func(), error) {
	panic(wire.Build(Sets))
}
