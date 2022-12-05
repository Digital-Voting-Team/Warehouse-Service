package main

import (
	"github.com/Digital-Voting-Team/warehouse-service/internal/cli"
	"os"
)

func main() {
	_ = os.Setenv("KV_VIPER_FILE", "config.yaml")

	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}
