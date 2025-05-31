package ai

import "errors"

var (
	ErrProviderNotFound = errors.New("provider not found")
)

type Provider interface {
	Ask(input *ProviderInput) (*ProviderOutput, error)
}

type ProviderInput struct {
	Model        string `json:"model"`
	Instructions string `json:"instructions"`
	Input        string `json:"input"`
}

type ProviderOutput struct {
	Status string
	Text   string
}

type ProviderFactory interface {
	Create(id string, apiKey string) (Provider, error)
}

type DefaultProviderFactory struct{}

func NewDefaultProviderFactory() *DefaultProviderFactory {
	return &DefaultProviderFactory{}
}

func (p *DefaultProviderFactory) Create(id string, apiKey string) (Provider, error) {
	providers := map[string]Provider{
		"openai": NewOpenAI(apiKey),
	}
	provider, providerExists := providers[id]
	if !providerExists {
		return nil, ErrProviderNotFound
	}
	return provider, nil
}
