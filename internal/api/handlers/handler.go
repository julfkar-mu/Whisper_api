package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	service "github.com/Whisper_api/internal/api/services"
	"github.com/Whisper_api/pkg/model"

	"github.com/gorilla/mux"
)

// Handler holds references to services used in the API handlers.
type Handler struct {
	WhisperService service.StreamService
}

// NewHandler creates a new Handler instance with required dependencies.
func NewHandler(s service.StreamService) *Handler {
	return &Handler{
		WhisperService: s,
	}
}

// RegisterRoutes registers API routes with the given router.
func (h *Handler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/api/transcribe", h.Transcribe).Methods("POST")
}

// Transcribe handles POST /api/transcribe to stream transcription output.
func (h *Handler) Transcribe(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form, limit max upload size (e.g., 50MB)
	if err := r.ParseMultipartForm(50 << 20); err != nil {
		http.Error(w, "Invalid multipart form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Missing file form field", http.StatusBadRequest)
		return
	}
	defer file.Close()

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Minute)
	defer cancel()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	// Pass the multipart file reader and filename to service which calls whisper client
	stream, errs := h.WhisperService.TranscribeStream(ctx, file, header.Filename)

	encoder := json.NewEncoder(w)

	for {
		select {
		case <-ctx.Done():
			return
		case err := <-errs:
			if err != nil {
				log.Printf("Error from transcription stream: %v", err)
				return
			}
		case text, ok := <-stream:
			if !ok {
				encoder.Encode(model.Transcription{Text: "", Done: true})
				flusher.Flush()
				return
			}
			encoder.Encode(model.Transcription{Text: text, Done: false})
			flusher.Flush()
		}
	}
}
