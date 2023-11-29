package httpchi

import (
	"github.com/rs/zerolog"
	"github.com/vlasashk/todo-manager/config"
	"github.com/vlasashk/todo-manager/internal/models/todo"
	"net/http"
)

type Service struct {
	DB todo.Repo
}

func NewService(db todo.Repo) Service {
	return Service{
		DB: db,
	}
}

func Run(service Service, logger zerolog.Logger, cfg config.AppCfg) {
	r := NewRouter(service, logger)
	logger.Info().Str("address", cfg.Host+":"+cfg.Port).Msg("starting listening")
	logger.Fatal().Err(http.ListenAndServe(cfg.Host+":"+cfg.Port, r))
}
