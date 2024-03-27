package main

import (
	"context"
	"log"

	"github.com/spf13/cobra"
)

func main() {
	ctx := context.Background()
	if err := rootCmd(ctx).Execute(); err != nil {
		log.Fatal(err)
	}
}

func rootCmd(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use: "cli",
	}

	cmd.AddCommand(
		serverCmd(ctx),
		userSubCmd(ctx),
		outputCmd(ctx),
	)

	return cmd
}
