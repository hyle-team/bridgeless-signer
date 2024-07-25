package main

import (
	"os"

	"github.com/hyle-team/bridgeless-signer/internal/cli"
)

func main() {
	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}
