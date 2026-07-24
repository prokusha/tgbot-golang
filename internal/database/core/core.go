package DB_CORE

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	ENV "github.com/prokusha/tgbot-golang/internal/env"
)

func DatabaseInit() (*pgxpool.Pool, error) {
	url := fmt.Sprintf("postgres://%s:%s@localhost:5432/%s", ENV.DB.User, ENV.DB.Password, ENV.DB.Database)
	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		return nil, err
	}
	return pool, nil
}
