package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseDSN string
	ServerPort  string
	JWTSecret   string
	BaseURL string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found")
	}

	dbHost := GetEnv("DB_HOST", "localhost")
	dbPort := GetEnv("DB_PORT", "5432")
	dbUser := GetEnv("DB_USER", "postgres")
	dbPasswd := GetEnv("DB_PASSWORD", "secret")
	dbName := GetEnv("DB_NAME", "url_shortener_db")
	sslMode := GetEnv("DB_SSLMODE", "disable")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", dbHost, dbPort, dbUser, dbPasswd, dbName, sslMode)
	return &Config{
		DatabaseDSN: dsn,
		ServerPort: GetEnv("SERVER_PORT", ":8080"),
		JWTSecret:   GetEnv("JWT_SECRET", "zxc"),
		BaseURL: GetEnv("BASE_URL", "http://localhost:8080"),
	}
}

func GetEnv(key, fallback string) string {
	if envVar := os.Getenv(key); envVar != "" {
		return envVar
	} else {
		return fallback
	}
}
