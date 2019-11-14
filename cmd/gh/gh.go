package main

import (
	"os"

	"gh/cmd/gh/app"
)

// Entrypoint for gh command
func main() {
	if err := app.Run(nil); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
