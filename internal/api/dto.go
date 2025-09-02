package api

import "github.com/Whisper_api/pkg/model"

type TranscriptionRequest struct {
	AudioURL string `json:"audio_url"`
}

// Alias or embed the shared model for response
type TranscriptionResponse = model.Transcription
