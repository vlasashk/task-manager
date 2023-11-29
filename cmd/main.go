package main

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/vlasashk/todo-manager/config"
	"github.com/vlasashk/todo-manager/internal/adapters/pgrepo"
	"github.com/vlasashk/todo-manager/internal/models/logger"
	"github.com/vlasashk/todo-manager/internal/ports/httpchi"
)

func main() {
	log := logger.NewLogger(zerolog.InfoLevel)
	log.Info().Msg("Logger created")
	cfg, err := config.ParseConfigValues()
	if err != nil {
		log.Fatal().Err(err).Msg("config parse fail")
	}
	log.Info().Msg("config parsing success")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	storage, err := pgrepo.NewTasksRepo(ctx, cfg.Postgres)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	log.Info().Msg("db connection success")
	service := httpchi.NewService(storage)
	httpchi.Run(service, log, cfg.App)
}
