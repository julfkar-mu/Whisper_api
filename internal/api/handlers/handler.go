package api

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

// TranscriptionRequest defines input payload for transcription (moved to dto or reuse).
type TranscriptionRequest struct {
	AudioURL string `json:"audio_url"` // URL of the audio file to transcribe
}

// Transcribe handles POST /api/transcribe to stream transcription output.
func (h *Handler) Transcribe(w http.ResponseWriter, r *http.Request) {
	// Decode request JSON
	var req TranscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Set a timeout context for the transcription streaming operation
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Minute)
	defer cancel()

	// Set response headers and status for streaming JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Get http.Flusher to flush responses as they stream
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	// Get transcription streaming channels from service
	stream, errs := h.WhisperService.TranscribeStream(ctx, req.AudioURL)

	encoder := json.NewEncoder(w)

	// Stream transcription results incrementally
	for {
		select {
		case <-ctx.Done():
			// Context timeout or client cancel
			return
		case err := <-errs:
			if err != nil {
				log.Printf("Error from transcription stream: %v", err)
				return
			}
		case text, ok := <-stream:
			if !ok {
				// Stream closed, send final done response
				encoder.Encode(model.Transcription{Text: "", Done: true})
				flusher.Flush()
				return
			}
			// Send partial transcription chunk
			encoder.Encode(model.Transcription{Text: text, Done: false})
			flusher.Flush()
		}
	}
}
