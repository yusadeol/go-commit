package ai

import "errors"

var (
	ErrProviderNotFound = errors.New("provider not found")
)

type Provider interface {
	Ask(input ProviderInput) (*ProviderOutput, error)
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

type ProviderFactory struct{}

func NewProviderFactory() *ProviderFactory {
	return &ProviderFactory{}
}

func (p *ProviderFactory) Create(providerName string) (Provider, error) {
	providers := map[string]Provider{
		"openai": NewOpenAI(),
	}
	provider, providerExists := providers[providerName]
	if !providerExists {
		return nil, ErrProviderNotFound
	}
	return provider, nil
}
