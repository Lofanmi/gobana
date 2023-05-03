package main

import (
	"time"
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
