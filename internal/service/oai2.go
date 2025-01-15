package service

import (
	"context"
	adapterApi "github.com/jiu-u/oai-adapter/api"
	"io"
	"net/http"
)

func (s *oaiService) ChatCompletions(ctx context.Context, req *adapterApi.ChatRequest) (io.ReadCloser, http.Header, error) {
	return s.RelayRequest(ctx, req, req.Model, RelayChat)
}

func (s *oaiService) ChatCompletionsByBytes(ctx context.Context, req []byte, modelId string) (io.ReadCloser, http.Header, error) {
	return s.RelayRequest(ctx, req, modelId, RelayChatByBytes)
}

func (s *oaiService) Completions(ctx context.Context, req *adapterApi.CompletionsRequest) (io.ReadCloser, http.Header, error) {
	return s.RelayRequest(ctx, req, req.Model, RelayCompletion)
}

func (s *oaiService) CompletionsByBytes(ctx context.Context, req []byte, modelId string) (io.ReadCloser, http.Header, error) {
	return s.RelayRequest(ctx, req, modelId, RelayCompletionByBytes)
}

func (s *oaiService) Embeddings(ctx context.Context, req *adapterApi.EmbeddingRequest) (io.ReadCloser, http.Header, error) {
	return s.RelayRequest(ctx, req, req.Model, RelayEmbedding)
}

func (s *oaiService) EmbeddingsByBytes(ctx context.Context, req []byte, modelId string) (io.ReadCloser, http.Header, error) {
	return s.RelayRequest(ctx, req, modelId, RelayEmbeddingByBytes)
}

func (s *oaiService) CreateSpeech(ctx context.Context, req *adapterApi.SpeechRequest) (io.ReadCloser, http.Header, error) {
	return s.RelayRequest(ctx, req, req.Model, RelaySpeech)
}

func (s *oaiService) CreateSpeechByBytes(ctx context.Context, req []byte, modelId string) (io.ReadCloser, http.Header, error) {
	return s.RelayRequest(ctx, req, modelId, RelaySpeechByBytes)
}

func (s *oaiService) Transcriptions(ctx context.Context, req *adapterApi.TranscriptionRequest) (io.ReadCloser, http.Header, error) {
	return s.RelayRequest(ctx, req, req.Model, RelayTranscriptions)
}

func (s *oaiService) Translations(ctx context.Context, req *adapterApi.TranslationRequest) (io.ReadCloser, http.Header, error) {
	return s.RelayRequest(ctx, req, req.Model, RelayTranslations)
}

func (s *oaiService) CreateImage(ctx context.Context, req *adapterApi.CreateImageRequest) (io.ReadCloser, http.Header, error) {
	return s.RelayRequest(ctx, req, req.Model, RelayImage)
}

func (s *oaiService) CreateImageByBytes(ctx context.Context, req []byte, modelId string) (io.ReadCloser, http.Header, error) {
	return s.RelayRequest(ctx, req, modelId, RelayImageByBytes)
}

func (s *oaiService) CreateImageEdit(ctx context.Context, req *adapterApi.EditImageRequest) (io.ReadCloser, http.Header, error) {
	return s.RelayRequest(ctx, req, req.Model, RelayImageEdit)
}

func (s *oaiService) ImageVariations(ctx context.Context, req *adapterApi.CreateImageVariationRequest) (io.ReadCloser, http.Header, error) {
	return s.RelayRequest(ctx, req, req.Model, RelayImageVariations)
}
