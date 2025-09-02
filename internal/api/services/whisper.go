package service

import (
	"context"
	"strings"

	"github.com/Whisper_api/internal/client"
)

// StreamService handles transcription streaming logic
type StreamService interface {
	TranscribeStream(ctx context.Context, audioURL string) (<-chan string, <-chan error)
}

type streamService struct {
	whisperClient client.WhisperClient
}

func NewStreamService(wc client.WhisperClient) StreamService {
	return &streamService{whisperClient: wc}
}

func (s *streamService) TranscribeStream(ctx context.Context, audioURL string) (<-chan string, <-chan error) {
	rawStream, errStream := s.whisperClient.StreamTranscription(ctx, audioURL)

	processedStream := make(chan string)
	processedErrors := make(chan error, 1)

	go func() {
		defer close(processedStream)
		defer close(processedErrors)

		for {
			select {
			case <-ctx.Done():
				processedErrors <- ctx.Err()
				return
			case text, ok := <-rawStream:
				if !ok {
					return
				}
				// Example process: clean or parse raw text
				processed := strings.TrimSpace(text)
				if processed != "" {
					processedStream <- processed
				}
			case err, ok := <-errStream:
				if ok && err != nil {
					processedErrors <- err
					return
				}
			}
		}
	}()
	return processedStream, processedErrors
}
