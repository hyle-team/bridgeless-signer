package service

import (
	"github.com/hyle-team/bridgeless-signer/internal/config"
	"gitlab.com/distributed_lab/logan/v3"
)

type service struct {
	log *logan.Entry
	cfg config.Config
}

func (s *service) run() error {
	s.log.Info("Service started")
	return nil
}

func newService(cfg config.Config) *service {
	return &service{
		log: cfg.Log(),
		cfg: cfg,
	}
}

func Run(cfg config.Config) {
	if err := newService(cfg).run(); err != nil {
		panic(err)
	}
}
