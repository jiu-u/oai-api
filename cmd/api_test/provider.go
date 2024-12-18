package main

import (
	"errors"
	"github.com/jiu-u/oai-api/pkg/adapter/provider"
)

func NewProvider(config provider.Config) (provider.Provider, error) {
	switch config.Type {
	case "openai":
		return provider.NewOpenAIProvider(config), nil
	default:
		return nil, errors.New("invalid provider type")
	}
}
