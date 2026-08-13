package main

import (
"context"
"fmt"
"log"

"github.com/your-org/hrms-backend/internal/config"
"github.com/your-org/hrms-backend/internal/db"
)

func main() {
cfg := config.Load()

log.Println("[INFO] Running migrations...")
if err := db.RunMigrations(cfg); err != nil {
log.Fatalf("[ERROR] Migration failed: %v", err)
}
log.Println("[SUCCESS] All database migrations applied successfully!")

// Verify tables
pgPool, err := db.NewPostgres(cfg)
if err != nil {
log.Fatalf("[ERROR] Failed to connect to DB: %v", err)
}
defer pgPool.Close()

rows, err := pgPool.Query(context.Background(),
"SELECT table_name FROM information_schema.tables WHERE table_schema='public' ORDER BY table_name;",
)
if err != nil {
log.Fatalf("[ERROR] Failed to query tables: %v", err)
}
defer rows.Close()

fmt.Println("\n=== CREATED TABLES IN DATABASE ===")
count := 0
for rows.Next() {
var name string
if err := rows.Scan(&name); err == nil {
count++
fmt.Printf(" [%02d] %s\n", count, name)
}
}
fmt.Printf("Total tables: %d\n", count)
}