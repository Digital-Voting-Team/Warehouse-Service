package main

import (
	"os"
	"warehouse-service/internal/cli"
)

func main() {
	_ = os.Setenv("KV_VIPER_FILE", "config.yaml")

	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}
