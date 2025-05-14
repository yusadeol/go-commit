package command

import (
	"github.com/yusadeol/go-commit/internal/application/usecase"
	"github.com/yusadeol/go-commit/internal/infrastructure/service/ai"
	"github.com/yusadeol/go-commit/internal/interface/cli"
)

type GenerateCommit struct{}

func NewGenerateCommit() *GenerateCommit {
	return &GenerateCommit{}
}

func (g GenerateCommit) GetName() string {
	return "generate"
}

func (g GenerateCommit) GetArguments() []cli.Argument {
	return []cli.Argument{
		{Name: "diff", Description: "Git diff", Required: false},
	}
}

func (g GenerateCommit) GetOptions() []cli.Option {
	return []cli.Option{
		{Name: "provider", Flag: "p", Description: "AI Provider", Default: "openai"},
		{Name: "language", Flag: "l", Description: "Language", Default: "en_US"},
	}
}

func (g GenerateCommit) Execute(input *cli.CommandInput) (*cli.ExecutionResult, error) {
	executionResult := cli.NewExecutionResult()
	providerFactory := ai.NewProviderFactory()
	aiProvider, err := providerFactory.Create(input.Options["provider"].Value)
	if err != nil {
		return nil, err
	}
	generateCommitInput := usecase.NewGenerateCommitInput(
		aiProvider,
		input.Options["language"].Value,
		input.Arguments["diff"].Value,
	)
	generateCommit := usecase.NewGenerateCommit()
	output, err := generateCommit.Execute(generateCommitInput)
	if err != nil {
		return nil, err
	}
	executionResult.Message = output.Commit
	return executionResult, nil
}
