package main

import (
	"github.com/mio256/wplus-server/pkg/ui"
)

func main() {
	r := ui.SetupRouter()
	r.Run(":8080")
}
