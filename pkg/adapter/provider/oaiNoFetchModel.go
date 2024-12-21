package provider

import "context"

type OaiNoFetchModelProvider struct {
	*OpenAIProvider
	models []string
}

func NewOaiNoFetchModelProvider(config Config, models []string) *OaiNoFetchModelProvider {
	return &OaiNoFetchModelProvider{
		OpenAIProvider: NewOpenAIProvider(config),
		models:         models,
	}
}

func (o *OaiNoFetchModelProvider) Models(ctx context.Context) ([]string, error) {
	return o.models, nil
}
