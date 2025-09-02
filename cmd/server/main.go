package main

import (
	"log"
	"net/http"
	"os"

	"whisperapi/internal/api"
	"whisperapi/internal/client"
	"whisperapi/internal/service"

	"github.com/gorilla/mux"
)

func main() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	whisperClient := client.NewOpenAIWhisperClient(apiKey)
	streamService := service.NewStreamService(whisperClient)
	handler := api.NewHandler(streamService)

	r := mux.NewRouter()
	api.RegisterRoutes(r, handler) // Use the routes.go to register routes

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
