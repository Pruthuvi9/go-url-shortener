package config

import "os"

type Config struct {
	Port       string
	DBDSN      string
	ValkeyAddr string
}

func Load() Config {
	return Config{
		Port:       getEnv("PORT", "3333"),
		DBDSN:      getEnv("DATABASE_URL", "postgres://localhost/urlshortner?sslmode=disable"),
		ValkeyAddr: getEnv("VALKEY_ADDR", "localhost:6379"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
