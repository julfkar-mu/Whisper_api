package client

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
)

// WhisperClient defines interface for OpenAI Whisper streaming client
type WhisperClient interface {
	StreamTranscription(ctx context.Context, audioURL string) (<-chan string, <-chan error)
}

type openAIWhisperClient struct {
	apiKey string
	apiURL string
}

func NewOpenAIWhisperClient(apiKey string) WhisperClient {
	return &openAIWhisperClient{
		apiKey: apiKey,
		apiURL: "https://api.openai.com/v1/audio/transcriptions/stream",
	}
}

func (c *openAIWhisperClient) StreamTranscription(ctx context.Context, audioURL string) (<-chan string, <-chan error) {
	texts := make(chan string)
	errs := make(chan error, 1)

	go func() {
		defer close(texts)
		defer close(errs)

		req, err := http.NewRequestWithContext(ctx, "POST", c.apiURL, nil) // Simplified, attach audioURL as needed
		if err != nil {
			errs <- err
			return
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
		// Add necessary headers and request body with audioURL or audio stream

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			errs <- err
			return
		}
		defer resp.Body.Close()

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				errs <- ctx.Err()
				return
			case texts <- scanner.Text():
			}
		}
		if err := scanner.Err(); err != nil {
			errs <- err
		}
	}()

	return texts, errs
}
