package main

import (
	"context"

	"github.com/maddiesch/marble/cmd/marble/cmd"
)

func main() {
	cmd.Execute(context.Background())
}
