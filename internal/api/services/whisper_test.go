package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// mockWhisperClient mocks client.WhisperClient interface
type mockWhisperClient struct {
	texts []string
	err   error
}

func (m *mockWhisperClient) StreamTranscription(ctx context.Context, audioURL string) (<-chan string, <-chan error) {
	textCh := make(chan string)
	errCh := make(chan error, 1)

	go func() {
		defer close(textCh)
		defer close(errCh)
		for _, t := range m.texts {
			textCh <- t
		}
		if m.err != nil {
			errCh <- m.err
		}
	}()
	return textCh, errCh
}

func TestTranscribeStream(t *testing.T) {
	client := &mockWhisperClient{
		texts: []string{"hello", "world"},
		err:   nil,
	}
	service := NewStreamService(client)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	textCh, errCh := service.TranscribeStream(ctx, "dummy-url")

	var results []string
	var err error
	done := false

	for !done {
		select {
		case text, ok := <-textCh:
			if !ok {
				done = true
				break
			}
			results = append(results, text)
		case err = <-errCh:
			done = true
		}
	}

	assert.NoError(t, err)
	assert.Equal(t, []string{"hello", "world"}, results)
}
