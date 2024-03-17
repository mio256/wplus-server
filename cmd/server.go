package main

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/mio256/wplus-server/pkg/ui"
	"github.com/spf13/cobra"
)

func serverCmd(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use: "server",
	}
	cmd.AddCommand(
		runLocalCmd(ctx),
		runNetworkCmd(ctx),
	)
	return cmd
}

func runLocalCmd(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use: "local",
		RunE: func(cmd *cobra.Command, args []string) error {
			r := ui.SetupRouter()

			if err := r.Run(":8080"); err != nil {
				panic(err)
			}

			return nil
		},
	}
	return cmd
}

func runNetworkCmd(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use: "network",
		RunE: func(cmd *cobra.Command, args []string) error {
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

			return nil
		},
	}
	return cmd
}
