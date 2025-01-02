package service

import (
	"context"
	adapterV1 "github.com/jiu-u/oai-api/pkg/adapter/api/v1"
	"io"
	"net/http"
)

func (s *oaiService) ChatCompletions(ctx context.Context, req *adapterV1.ChatCompletionRequest) (io.ReadCloser, http.Header, error) {
	return s.RelayRequest(ctx, req, req.Model, RelayChat)
}

func (s *oaiService) ChatCompletionsByBytes(ctx context.Context, req []byte, modelId string) (io.ReadCloser, http.Header, error) {
	return s.RelayRequest(ctx, req, modelId, RelayChatByBytes)
}

func (s *oaiService) Completions(ctx context.Context, req *adapterV1.CompletionsRequest) (io.ReadCloser, http.Header, error) {
	return s.RelayRequest(ctx, req, req.Model, RelayCompletion)
}

func (s *oaiService) CompletionsByBytes(ctx context.Context, req []byte, modelId string) (io.ReadCloser, http.Header, error) {
	return s.RelayRequest(ctx, req, modelId, RelayCompletionByBytes)
}

func (s *oaiService) Embeddings(ctx context.Context, req *adapterV1.EmbeddingRequest) (io.ReadCloser, http.Header, error) {
	return s.RelayRequest(ctx, req, req.Model, RelayEmbedding)
}

func (s *oaiService) EmbeddingsByBytes(ctx context.Context, req []byte, modelId string) (io.ReadCloser, http.Header, error) {
	return s.RelayRequest(ctx, req, modelId, RelayEmbeddingByBytes)
}

func (s *oaiService) CreateSpeech(ctx context.Context, req *adapterV1.SpeechRequest) (io.ReadCloser, http.Header, error) {
	return s.RelayRequest(ctx, req, req.Model, RelaySpeech)
}

func (s *oaiService) CreateSpeechByBytes(ctx context.Context, req []byte, modelId string) (io.ReadCloser, http.Header, error) {
	return s.RelayRequest(ctx, req, modelId, RelaySpeechByBytes)
}

func (s *oaiService) Transcriptions(ctx context.Context, req *adapterV1.TranscriptionRequest) (io.ReadCloser, http.Header, error) {
	return s.RelayRequest(ctx, req, req.Model, RelayTranscriptions)
}

func (s *oaiService) Translations(ctx context.Context, req *adapterV1.TranslationRequest) (io.ReadCloser, http.Header, error) {
	return s.RelayRequest(ctx, req, req.Model, RelayTranslations)
}

func (s *oaiService) CreateImage(ctx context.Context, req *adapterV1.CreateImageRequest) (io.ReadCloser, http.Header, error) {
	return s.RelayRequest(ctx, req, req.Model, RelayImage)
}

func (s *oaiService) CreateImageByBytes(ctx context.Context, req []byte, modelId string) (io.ReadCloser, http.Header, error) {
	return s.RelayRequest(ctx, req, modelId, RelayImageByBytes)
}

func (s *oaiService) CreateImageEdit(ctx context.Context, req *adapterV1.EditImageRequest) (io.ReadCloser, http.Header, error) {
	return s.RelayRequest(ctx, req, req.Model, RelayImageEdit)
}

func (s *oaiService) ImageVariations(ctx context.Context, req *adapterV1.CreateImageVariationRequest) (io.ReadCloser, http.Header, error) {
	return s.RelayRequest(ctx, req, req.Model, RelayImageVariations)
}
