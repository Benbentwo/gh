package main

import (
	"os"

	"github.com/Benbentwo/gh/cmd/gh/app"
)

// Entrypoint for gh command
func main() {
	if err := app.Run(nil); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
