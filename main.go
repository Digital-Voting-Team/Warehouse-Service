package main

import (
	"github.com/Digital-Voting-Team/warehouse-service/internal/cli"
	"os"
)

func main() {
	//_ = os.Setenv("KV_VIPER_FILE", "config.yaml")
	//_ = os.Setenv("DB_URL", "oracle://WAREHOUSE:WAREHOUSE@warehouse-db:1521/XEPDB1")

	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}
