package main

import (
	"github/ptsypyshev/go-observability/lesson03/app"
	"log"
)

func main() {
	a := app.App{}
	if closer, err := a.Init(); err != nil {
		log.Fatal(err)
	} else {
		defer closer.Close()
	}

	if err := a.Serve(); err != nil {
		log.Fatal(err)
	}
}
