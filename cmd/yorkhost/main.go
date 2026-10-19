// Command yorkhost : CLI officiel YorkHost.
package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/YorkHost-fr/yorkhost-cli/internal/cli"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	code := cli.Execute(ctx, cli.NewApp(), os.Args[1:])
	stop()
	os.Exit(code)
}
