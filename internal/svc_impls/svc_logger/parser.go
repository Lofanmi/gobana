package svc_logger

import (
	"context"

	"github.com/Lofanmi/gobana/service"
)

type LogParser interface {
	Do(ctx context.Context, raw string) service.LogItem
}

type logParser struct{}

func (logParser) Do(ctx context.Context, raw string) service.LogItem {
	// TODO implement me
	panic("implement me")
}
