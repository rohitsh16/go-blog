package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/rohitsh16/go-blog/backend/server"
)

func main() {
	srv, err := server.NewServer()
	if err != nil {
		log.Fatalf("Server initialization failed: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := srv.Start(); err != nil && err.Error() != "http: Server closed" {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Termination signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Shutdown error: %v", err)
	}

	log.Println("Server exited gracefully")
}
