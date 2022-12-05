package cli

import (
	_ "github.com/sijms/go-ora/v2"

	"github.com/Digital-Voting-Team/warehouse-service/internal/config"
	"github.com/Digital-Voting-Team/warehouse-service/internal/service"
	"github.com/jmoiron/sqlx"

	"github.com/alecthomas/kingpin"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3"
)

func Run(args []string) bool {
	log := logan.New()

	defer func() {
		if rvr := recover(); rvr != nil {
			log.WithRecover(rvr).Error("app panicked")
		}
	}()

	cfg := config.New(kv.MustFromEnv())
	log = cfg.Log()

	app := kingpin.New("warehouse-service", "")

	runCmd := app.Command("run", "run command")
	serviceCmd := runCmd.Command("service", "run service") // you can insert custom help

	// custom commands go here...

	db, err := sqlx.Connect("oracle", "oracle://WAREHOUSE:WAREHOUSE@warehouse-db:1521/XEPDB1")

	if err != nil {
		log.Panic(err)
	}

	cmd, err := app.Parse(args[1:])
	if err != nil {
		log.WithError(err).Error("failed to parse arguments")
		return false
	}

	switch cmd {
	case serviceCmd.FullCommand():
		service.Run(cfg, db)
	// handle any custom commands here in the same way
	default:
		log.Errorf("unknown command %s", cmd)
		return false
	}

	return true
}
