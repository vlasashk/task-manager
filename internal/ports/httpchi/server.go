package httpchi

import (
	"github.com/rs/zerolog"
	"github.com/vlasashk/task-manager/config"
	"github.com/vlasashk/task-manager/internal/models/tasktodo"
	"net/http"
)

type Service struct {
	DB tasktodo.Repo
}

func NewService(db tasktodo.Repo) Service {
	return Service{
		DB: db,
	}
}

func Run(service Service, logger zerolog.Logger, cfg config.AppCfg) {
	r := NewRouter(service, logger)
	logger.Info().Str("address", cfg.Host+":"+cfg.Port).Msg("starting listening")
	logger.Fatal().Err(http.ListenAndServe(cfg.Host+":"+cfg.Port, r))
}
