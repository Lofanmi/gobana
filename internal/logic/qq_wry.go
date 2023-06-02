package logic

import (
	"context"
	"fmt"
)

type QQWry interface {
	Find(ctx context.Context, ip string) (location fmt.Stringer, err error)
}
