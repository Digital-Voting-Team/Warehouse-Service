package service

import (
	"github.com/jmoiron/sqlx"
	"net"
	"net/http"

	"warehouse-service/internal/config"

	"gitlab.com/distributed_lab/kit/copus/types"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type service struct {
	log       *logan.Entry
	copus     types.Copus
	listener  net.Listener
	endpoints *config.EndpointsConfig
}

func (s *service) run(db *sqlx.DB) error {
	s.log.Info("Service started")
	r := s.router(db)

	if err := s.copus.RegisterChi(r); err != nil {
		return errors.Wrap(err, "cop failed")
	}

	return http.Serve(s.listener, r)
}

func newService(cfg config.Config) *service {
	return &service{
		log:       cfg.Log(),
		copus:     cfg.Copus(),
		listener:  cfg.Listener(),
		endpoints: cfg.EndpointsConfig(),
	}
}

func Run(cfg config.Config, db *sqlx.DB) {
	if err := newService(cfg).run(db); err != nil {
		panic(err)
	}
}
