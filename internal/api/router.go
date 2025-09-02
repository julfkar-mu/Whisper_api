package api

import (
	"github.com/gorilla/mux"
)

// RegisterRoutes registers all API routes for the application
func RegisterRoutes(r *mux.Router, handler *Handler) {
	// Transcription streaming endpoint
	r.HandleFunc("/api/transcribe", handler.Transcribe).Methods("POST")

	// Additional routes can be added here in the future
}
