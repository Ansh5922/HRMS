package db

import (
"database/sql"
"fmt"
"log"

"github.com/golang-migrate/migrate/v4"
"github.com/golang-migrate/migrate/v4/database/postgres"
_ "github.com/golang-migrate/migrate/v4/source/file"
_ "github.com/lib/pq"
"github.com/your-org/hrms-backend/internal/config"
)

// RunMigrations applies all pending UP migrations from the migrations/ directory.
func RunMigrations(cfg *config.Config) error {
dsn := fmt.Sprintf(
"postgres://%s:%s@%s:%s/%s?sslmode=%s",
cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBSSLMode,
)

db, err := sql.Open("postgres", dsn)
if err != nil {
return fmt.Errorf("open db for migrations: %w", err)
}
defer db.Close()

driver, err := postgres.WithInstance(db, &postgres.Config{})
if err != nil {
return fmt.Errorf("create migration driver: %w", err)
}

m, err := migrate.NewWithDatabaseInstance(
"file://migrations",
"postgres",
driver,
)
if err != nil {
return fmt.Errorf("create migrate instance: %w", err)
}

if err := m.Up(); err != nil && err != migrate.ErrNoChange {
return fmt.Errorf("run migrations: %w", err)
}

version, dirty, _ := m.Version()
log.Printf("[DB] Migrations applied — version: %d, dirty: %v", version, dirty)
return nil
}
