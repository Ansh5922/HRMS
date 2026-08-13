package db

import (
"context"
"fmt"

"github.com/jackc/pgx/v5/pgxpool"
"github.com/your-org/hrms-backend/internal/config"
)

func NewPostgres(cfg *config.Config) (*pgxpool.Pool, error) {
dsn := fmt.Sprintf(
"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s pool_max_conns=25",
cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBUser, cfg.DBPassword, cfg.DBSSLMode,
)
pool, err := pgxpool.New(context.Background(), dsn)
if err != nil {
return nil, err
}
if err := pool.Ping(context.Background()); err != nil {
return nil, err
}
return pool, nil
}
