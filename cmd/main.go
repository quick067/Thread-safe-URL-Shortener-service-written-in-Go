package main

import (
	cfg "URL-shortener/internal/config"
	handlers "URL-shortener/internal/handlers"
	"URL-shortener/internal/logger"
	"URL-shortener/internal/middleware"
	serv "URL-shortener/internal/server"
	service "URL-shortener/internal/service"
	"URL-shortener/internal/store"

	"log"
	"math/rand"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)


func main() {
	logger.SetupLogger()

	rand.Seed(time.Now().UnixNano())
	serverCfg := cfg.Load()
	
	storage, err := store.InitDB(serverCfg.DatabaseDSN)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer storage.Close()
	storage.CreateTables()

	authService := &service.AuthService{
		Storage: storage,
		JWTSecret: serverCfg.JWTSecret,
	}
	linkService := &service.LinkService{
		Storage: storage,
		BaseURL: serverCfg.BaseURL,
	}

	handler := &handlers.Handler{
		AuthService: authService,
		LinkService: linkService,
	}

	handler.RegisterRoutes()
	mux := http.DefaultServeMux

	rl := &middleware.RateLimiter{
		LimitMap: make(map[string]*rate.Limiter),
	}
	limitRouter := rl.RateLimitMiddleware(mux)
	loggerRouter := middleware.LogMiddleware(limitRouter)
	finalHandler := middleware.MetricsMiddleware(loggerRouter)

	serv.StartServer(serverCfg.ServerPort, finalHandler, nil)
}
