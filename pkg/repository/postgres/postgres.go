package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Hymiside/hezzl-api/pkg/models"
	_ "github.com/lib/pq"
)

func NewPostgresDB(ctx context.Context, cfg models.ConfigPostgres) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("error to connection postgres: %v", err)
	}
	go func(ctx context.Context) {
		<-ctx.Done()
		db.Close()
	}(ctx)

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("connection test error: %w", err)
	}

	return db, nil
}