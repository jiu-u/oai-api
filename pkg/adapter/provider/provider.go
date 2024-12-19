package provider

import (
	"context"
	"errors"
	"github.com/jiu-u/oai-api/pkg/adapter/api/v1"
	"io"
	"net/http"
)

type Config struct {
	Type     string
	EndPoint string
	APIKey   string
}

type Provider interface {
	ChatCompletions(ctx context.Context, req *v1.ChatCompletionRequest) (io.Reader, http.Header, error)
	ChatCompletionsByBytes(ctx context.Context, req []byte) (io.Reader, http.Header, error)
	Models(ctx context.Context) ([]string, error)
	Completions(ctx context.Context, req *v1.CompletionsRequest) (io.Reader, http.Header, error)
	CompletionsByBytes(ctx context.Context, req []byte) (io.Reader, http.Header, error)
	Embeddings(ctx context.Context, req *v1.EmbeddingRequest) (io.Reader, http.Header, error)
	EmbeddingsByBytes(ctx context.Context, req []byte) (io.Reader, http.Header, error)
	CreateSpeech(ctx context.Context, req *v1.SpeechRequest) (io.Reader, http.Header, error)
	CreateSpeechByBytes(ctx context.Context, req []byte) (io.Reader, http.Header, error)
	Transcriptions(ctx context.Context, req *v1.TranscriptionRequest) (io.Reader, http.Header, error)
	Translations(ctx context.Context, req *v1.TranslationRequest) (io.Reader, http.Header, error)
	CreateImage(ctx context.Context, req *v1.CreateImageRequest) (io.Reader, http.Header, error)
	CreateImageByBytes(ctx context.Context, req []byte) (io.Reader, http.Header, error)
	CreateImageEdit(ctx context.Context, req *v1.EditImageRequest) (io.Reader, http.Header, error)
	ImageVariations(ctx context.Context, req *v1.CreateImageVariationRequest) (io.Reader, http.Header, error)
}

func HandleUnSupportedError() (io.Reader, http.Header, error) {
	return nil, nil, errors.New("the feature is not supported")
}
