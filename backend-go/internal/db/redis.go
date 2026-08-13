package db

import (
"fmt"

"github.com/redis/go-redis/v9"
"github.com/your-org/hrms-backend/internal/config"
)

func NewRedis(cfg *config.Config) *redis.Client {
return redis.NewClient(&redis.Options{
Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
Password: cfg.RedisPassword,
DB:       0,
})
}
