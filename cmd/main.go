package main

import (
	"time"

	"github.com/Lofanmi/gobana/cmd/inject"
	"github.com/json-iterator/go/extra"
)

func main() {
	extra.RegisterFuzzyDecoders()
	time.Local, _ = time.LoadLocation("PRC")
	application, cleanup, err := inject.NewApplication()
	if err != nil {
		panic(err)
	}
	defer cleanup()
	application.Run()
}
