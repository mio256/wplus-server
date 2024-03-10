package main

import (
	"fmt"
	"net"
	"os"

	"github.com/mio256/wplus-server/pkg/ui"
)

func main() {
	r := ui.SetupRouter()
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				r.Run(fmt.Sprint(ipnet.IP.String(), ":8080"))
			}
		}
	}
}
