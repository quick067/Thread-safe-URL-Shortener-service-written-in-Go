package main

import (
	cfg "URL-shortener/internal/config"
	handlers "URL-shortener/internal/handlers"
	serv "URL-shortener/internal/server"
	"URL-shortener/internal/store"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	port := cfg.GetEnv("SERVER_PORT", ":8080")

	storage := store.New(cfg.GetEnv("FILENAME", "storage.json"))
	storage.LoadFromFile()

	handler := handlers.Handler{
		Storage: storage,
	}

	handler.RegisterRoutes()

	serv.StartServer(port, storage.SaveToFile)
}
