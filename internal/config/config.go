package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	RedisURL    string
	BaseURL     string
	Secret      string
	CacheTTL    time.Duration
	AppEnv      string
}

var C Config

func init() {
	// грузим .env (если его нет, пропускаем)
	_ = godotenv.Load()

	ttl, err := strconv.Atoi(getEnv("CACHE_TTL_SECONDS", "300"))
	if err != nil {
		log.Fatalf("invalid CACHE_TTL_SECONDS: %v", err)
	}

	C = Config{
		DatabaseURL: getEnv("DATABASE_URL", ""),
		RedisURL:    getEnv("REDIS_URL", ""),
		BaseURL:     getEnv("BASE_URL", "http://localhost:3000"),
		Secret:      getEnv("SHORT_CODE_SECRET", "defaultsecret"),
		CacheTTL:    time.Duration(ttl) * time.Second,
		AppEnv:      getEnv("APP_ENV", "development"),
	}
}

func getEnv(key, def string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return def
}
