package config

import (
"os"
"github.com/joho/godotenv"
)

type Config struct {
AppPort       string
AppEnv        string
AppSecret     string
DBHost        string
DBPort        string
DBName        string
DBUser        string
DBPassword    string
DBSSLMode     string
RedisHost     string
RedisPort     string
RedisPassword string
JWTAccessExp  string
JWTRefreshExp string
StorageURL    string
StorageBucket string
AIServiceURL  string
AIServiceSecret string
}

func Load() *Config {
_ = godotenv.Load()
return &Config{
AppPort:         getEnv("APP_PORT", "8080"),
AppEnv:          getEnv("APP_ENV", "development"),
AppSecret:       getEnv("APP_SECRET", ""),
DBHost:          getEnv("DB_HOST", "localhost"),
DBPort:          getEnv("DB_PORT", "5432"),
DBName:          getEnv("DB_NAME", "hrms_dev"),
DBUser:          getEnv("DB_USER", "postgres"),
DBPassword:      getEnv("DB_PASSWORD", ""),
DBSSLMode:       getEnv("DB_SSLMODE", "disable"),
RedisHost:       getEnv("REDIS_HOST", "localhost"),
RedisPort:       getEnv("REDIS_PORT", "6379"),
RedisPassword:   getEnv("REDIS_PASSWORD", ""),
JWTAccessExp:    getEnv("JWT_ACCESS_EXPIRY", "15m"),
JWTRefreshExp:   getEnv("JWT_REFRESH_EXPIRY", "7d"),
StorageURL:      getEnv("STORAGE_ENDPOINT", ""),
StorageBucket:   getEnv("STORAGE_BUCKET", "hrms"),
AIServiceURL:    getEnv("AI_SERVICE_URL", "http://localhost:8001"),
AIServiceSecret: getEnv("AI_SERVICE_SECRET", ""),
}
}

func getEnv(key, fallback string) string {
if v := os.Getenv(key); v != "" {
return v
}
return fallback
}
