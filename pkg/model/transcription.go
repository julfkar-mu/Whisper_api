package model

// Transcription represents a single chunk of transcription text from Whisper.
type Transcription struct {
	Text string `json:"text"`
	Done bool   `json:"done"` // Indicates if the transcription is complete
}
