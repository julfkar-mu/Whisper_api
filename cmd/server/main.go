package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Whisper_api/internal/api"
	apihandler "github.com/Whisper_api/internal/api/handlers"
	service "github.com/Whisper_api/internal/api/services"
	"github.com/Whisper_api/internal/client"

	"github.com/gorilla/mux"
)

func main() {
	apiKey := os.Getenv("OPENAI_API_KEY")

	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable is required")
	}

	whisperClient := client.NewOpenAIWhisperClient(apiKey)
	streamService := service.NewStreamService(whisperClient)
	handler := apihandler.NewHandler(streamService)

	r := mux.NewRouter()
	api.RegisterRoutes(r, handler) // Use the routes.go to register routes

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Channel to listen for interrupt or termination signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Run server in goroutine so that it doesn't block
	go func() {
		log.Println("Starting server on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	// Wait for interrupt signal
	<-stop
	log.Println("Shutting down server gracefully...")

	// Create context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown server gracefully
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Println("Server stopped")
}
