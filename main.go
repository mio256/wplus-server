package main

import (
	"github.com/mio256/wplus-server/pkg/ui"
)

func main() {
	r := ui.SetupRouter()
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
