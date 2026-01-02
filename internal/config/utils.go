package config

import (
	"os"
)

func GetEnv(key, fallback string) string {
	if envVar := os.Getenv(key); envVar != "" {
		return envVar
	} else {
		return fallback
	}
}
