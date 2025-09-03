package api

import (
	handler "github.com/Whisper_api/internal/api/handlers"
	"github.com/gorilla/mux"
)

// RegisterRoutes registers all API routes for the application
func RegisterRoutes(r *mux.Router, handler *handler.Handler) {
	// Transcription streaming endpoint
	r.HandleFunc("/api/transcribe", handler.Transcribe).Methods("POST")

	// Additional routes can be added here in the future
}
