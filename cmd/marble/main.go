package main

import (
	"context"
	"fmt"
	"os"

	"github.com/maddiesch/marble/cmd/marble/cmd"
)

func main() {
	if err := cmd.Marble.ExecuteContext(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "[FATAL] %s\n", err.Error())
		os.Exit(1)
	}
}
