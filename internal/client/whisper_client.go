package client

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
)

// WhisperClient interface updated to take io.Reader (audio file stream)
type WhisperClient interface {
	StreamTranscription(ctx context.Context, audioFile io.Reader, fileName string) (<-chan string, <-chan error)
}

type openAIWhisperClient struct {
	apiKey string
	apiURL string
}

func NewOpenAIWhisperClient(apiKey string) WhisperClient {
	return &openAIWhisperClient{
		apiKey: apiKey,
		apiURL: "https://api.openai.com/v1/audio/transcriptions",
	}
}

func (c *openAIWhisperClient) StreamTranscription(ctx context.Context, audioFile io.Reader, fileName string) (<-chan string, <-chan error) {
	texts := make(chan string)
	errs := make(chan error, 1)

	log.Printf("Audio file: %s", fileName)

	go func() {
		defer close(texts)
		defer close(errs)

		pr, pw := io.Pipe()
		writer := multipart.NewWriter(pw)

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.apiURL, pr)
		if err != nil {
			errs <- fmt.Errorf("failed to create request: %w", err)
			return
		}

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
		req.Header.Set("Content-Type", writer.FormDataContentType())
		// Write multipart form asynchronously
		go func() {
			var closeErr error
			defer func() {
				// Ensure pipe is closed with an error (if any)
				writer.Close()
				pw.CloseWithError(closeErr)
			}()

			part, err := writer.CreateFormFile("file", fileName)
			if err != nil {
				closeErr = fmt.Errorf("failed to create form file field: %w", err)
				errs <- closeErr
				return
			}

			if written, err := io.Copy(part, audioFile); err != nil {
				closeErr = fmt.Errorf("failed to copy audio file data: %w", err)
				errs <- closeErr
				return
			} else {
				log.Printf("Written bytes to form file:", written)
			}

			if err := writer.WriteField("model", "gpt-4o-mini-transcribe"); err != nil {
				closeErr = fmt.Errorf("failed to write model field: %w", err)
				errs <- closeErr
				return
			}

			if err := writer.WriteField("stream", strconv.FormatBool(true)); err != nil {
				closeErr = fmt.Errorf("failed to write stream field: %w", err)
				errs <- closeErr
				return
			}

			log.Println("Finished writing multipart form")
		}()

		log.Println("Sending request to OpenAI Whisper API")
		resp, err := http.DefaultClient.Do(req)

		if err != nil {
			errs <- fmt.Errorf("request error: %w", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Printf("Error response status:%s", err)
			j, _ := io.ReadAll(resp.Body)
			log.Println("Error response body:", string(j))
			errs <- fmt.Errorf("unexpected status code: %d", resp.StatusCode)
			return
		}

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
			errs <- fmt.Errorf("error reading stream: %w", err)
		}
	}()

	return texts, errs
}
