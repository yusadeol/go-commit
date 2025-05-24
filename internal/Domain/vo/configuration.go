package vo

type Configuration struct {
	DefaultAIProvider string                `json:"default_ai_provider"`
	DefaultLanguage   string                `json:"default_language"`
	AIProviders       map[string]AIProvider `json:"ai_providers"`
	Languages         map[string]Language   `json:"languages"`
}

type AIProvider struct {
	ID           string         `json:"id"`
	APIKey       string         `json:"api_key"`
	Models       []string       `json:"models"`
	DefaultModel string         `json:"default_model"`
	Enabled      bool           `json:"enabled"`
	Settings     map[string]any `json:"settings"`
}

type Language struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	Enabled     bool   `json:"enabled"`
}
