package provider

type SiliconFlowProvider struct {
	*OpenAIProvider
}

func NewSiliconFlowProvider(config Config) *SiliconFlowProvider {
	return &SiliconFlowProvider{
		OpenAIProvider: NewOpenAIProvider(config),
	}
}
