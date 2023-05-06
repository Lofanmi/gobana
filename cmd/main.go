package main

import (
	"time"

	"github.com/Lofanmi/gobana/cmd/inject"
)

func main() {
	time.Local, _ = time.LoadLocation("PRC")
	application, cleanup, err := inject.NewApplication()
	if err != nil {
		panic(err)
	}
	defer cleanup()
	application.Run()
}
