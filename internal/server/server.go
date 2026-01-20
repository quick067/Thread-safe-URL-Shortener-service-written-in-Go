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

func StartServer(addr string, handler http.Handler, cleanFunc func() error) {
	server := http.Server{
		Addr:    addr,
		Handler: handler,
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
	log.Printf("Closing DB connection...")
	if cleanFunc != nil{
		if err := cleanFunc(); err != nil {
			log.Printf("error executing cleanFunc: %v", err)
		}
	}
	log.Println("Server stopped gracefully")
}
