package pgrepo

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vlasashk/todo-manager/config"
	"os"
	"time"
)

type Repo struct {
	DB  *pgxpool.Pool
	ctx context.Context
}

func NewTasksRepo(ctx context.Context, cfg config.PostgresCfg) (Repo, error) {
	url := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.NameDB)
	dbPool, err := pgxpool.New(ctx, url)
	if err != nil {
		return Repo{}, fmt.Errorf("unable to create connection pool: %v", err)
	}
	timeCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err = dbPool.Ping(timeCtx); err != nil {
		return Repo{}, fmt.Errorf("unable to ping connection pool: %v", err)
	}
	instance := Repo{dbPool, ctx}
	if err = instance.NewTable(timeCtx, cfg); err != nil {
		return Repo{}, err
	}
	return instance, nil
}

func (db Repo) NewTable(ctx context.Context, cfg config.PostgresCfg) error {
	err := db.DB.AcquireFunc(ctx, func(conn *pgxpool.Conn) error {
		query, err := os.ReadFile(cfg.InitFilePath)
		if err != nil {
			return fmt.Errorf("failed to read sql file: %v", err)
		}
		if _, err = conn.Exec(ctx, string(query)); err != nil {
			return fmt.Errorf("failed to init tables: %v", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("unable to acquire a database connection: %v", err)
	}
	return nil
}
