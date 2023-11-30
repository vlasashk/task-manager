//	@title			task-manager API
//	@version		1.0
//	@description	API for task manager

//	@host		localhost:9090
//	@BasePath	/api/

// @securityDefinitions.basic	BasicAuth
package main

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/vlasashk/task-manager/config"
	"github.com/vlasashk/task-manager/internal/adapters/pgrepo"
	"github.com/vlasashk/task-manager/internal/models/logger"
	"github.com/vlasashk/task-manager/internal/ports/httpchi"
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
