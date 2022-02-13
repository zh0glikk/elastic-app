package service

import (
	"context"
	"github.com/sirupsen/logrus"
	"net/http"

	"github.com/zh0glikk/elastic-app/internal/api"
	"github.com/zh0glikk/elastic-app/internal/config"
)

type Service struct {
	cfg *config.Config
	log *logrus.Entry
}

func NewService(cfg *config.Config) *Service {
	switch {
	case cfg == nil:
		panic("cfg is nil")
	}

	return &Service{
		cfg: cfg,
		log: logrus.NewEntry(logrus.New()),
	}
}

func (s *Service) Run(ctx context.Context) error {
	router, err := api.Router(s.log, *s.cfg)
	if err != nil {
		return err
	}

	return http.ListenAndServe(s.cfg.ListenerConfig.Addr, router)
}

