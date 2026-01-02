package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func StartServer(addr string, cleanFunc func() error) {
	server := http.Server{
		Addr:    addr,
		Handler: nil,
	}
	log.Printf("Server starting on port: %s ...", server.Addr)

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed && err != nil {
			log.Fatalf("error listening server: %v", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
	log.Printf("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("error shutting down the server: %v", err)
	}
	log.Printf("Saving file...")
	if err := cleanFunc(); err != nil {
		log.Printf("error saving file: %v", err)
	}
	log.Println("Server stopped gracefully")
}
