package main

import (
	"context"
	"log"

	"github.com/your-org/hrms-backend/internal/config"
	"github.com/your-org/hrms-backend/internal/db"
	"github.com/your-org/hrms-backend/internal/router"
)

func main() {
	// 1. Load config from .env
	cfg := config.Load()

	// 2. Run database migrations
	if err := db.RunMigrations(cfg); err != nil {
		log.Fatalf("[FATAL] Migration failed: %v", err)
	}

	// 3. Connect PostgreSQL pool
	pgPool, err := db.NewPostgres(cfg)
	if err != nil {
		log.Fatalf("[FATAL] PostgreSQL connection failed: %v", err)
	}
	defer pgPool.Close()
	log.Println("[DB] PostgreSQL connected")

	// 4. Connect Redis
	redisClient := db.NewRedis(cfg)
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Printf("[WARN] Redis connection warning: %v (continuing with DB fallback)", err)
	} else {
		log.Println("[DB] Redis connected")
	}

	// 5. Set up HTTP router and start server
	r := router.Setup(cfg, pgPool, redisClient)
	log.Printf("[SERVER] Listening on :%s (env: %s)", cfg.AppPort, cfg.AppEnv)
	log.Printf("[SWAGGER] Documentation available at http://localhost:%s/swagger/index.html", cfg.AppPort)

	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("[FATAL] Server error: %v", err)
	}
}
