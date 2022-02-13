package main

import (
	"os"

	"github.com/zh0glikk/elastic-app/internal/cli"
)

func main() {
	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}
