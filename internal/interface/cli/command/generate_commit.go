package command

import (
	"bytes"
	"errors"
	"github.com/yusadeol/go-commit/internal/application/usecase"
	"github.com/yusadeol/go-commit/internal/infrastructure/service/ai"
	"github.com/yusadeol/go-commit/internal/interface/cli"
	"os/exec"
)

type GenerateCommit struct {
	configuration *Configuration
}

func NewGenerateCommit(configuration *Configuration) *GenerateCommit {
	return &GenerateCommit{configuration: configuration}
}

type Configuration struct {
	DefaultAIProvider string                `json:"default_ai_provider"`
	DefaultLanguage   string                `json:"default_language"`
	AIProviders       map[string]AIProvider `json:"ai_providers"`
	Languages         map[string]Language   `json:"languages"`
}

type AIProvider struct {
	APIKey       string                 `json:"api_key"`
	Models       []string               `json:"models"`
	DefaultModel string                 `json:"default_model"`
	Enabled      bool                   `json:"enabled"`
	Settings     map[string]interface{} `json:"settings"`
}

type Language struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

func (g *GenerateCommit) GetName() string {
	return "generate"
}

func (g *GenerateCommit) GetArguments() []cli.Argument {
	return []cli.Argument{
		{Name: "diff", Description: "Git diff", Required: false},
	}
}

func (g *GenerateCommit) GetOptions() []cli.Option {
	return []cli.Option{
		{Name: "provider", Flag: "p", Description: "AI Provider", Default: "openai"},
		{Name: "language", Flag: "l", Description: "Language", Default: "en_US"},
	}
}

func (g *GenerateCommit) Execute(input *cli.CommandInput) (*cli.ExecutionResult, error) {
	executionResult := cli.NewExecutionResult()
	providerFactory := ai.NewProviderFactory()
	configurationAiProvider, configurationAiProviderExists := g.configuration.AIProviders[input.Options["provider"].Value]
	if !configurationAiProviderExists {
		return nil, errors.New("AI Provider configuration not found")
	}
	aiProvider, err := providerFactory.Create(input.Options["provider"].Value, configurationAiProvider.APIKey)
	if err != nil {
		return nil, err
	}
	diff := input.Arguments["diff"].Value
	if diff == "" {
		diff, err = g.GetGitDiff()
		if err != nil {
			return nil, err
		}
	}
	generateCommitInput := usecase.NewGenerateCommitInput(
		configurationAiProvider.DefaultModel,
		aiProvider,
		input.Options["language"].Value,
		diff,
	)
	generateCommit := usecase.NewGenerateCommit()
	output, err := generateCommit.Execute(generateCommitInput)
	if err != nil {
		return nil, err
	}
	executionResult.Message = output.Commit
	return executionResult, nil
}

func (g *GenerateCommit) GetGitDiff() (string, error) {
	var out bytes.Buffer
	cmd := exec.Command("git", "diff", "--staged")
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	diff := out.String()
	if diff == "" {
		return "", errors.New("no staged changes found")
	}
	return diff, nil
}
