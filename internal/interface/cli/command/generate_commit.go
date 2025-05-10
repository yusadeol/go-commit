package command

import (
	"fmt"
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
	return []cli.Argument{}
}

func (g GenerateCommit) GetOptions() []cli.Option {
	return []cli.Option{
		{Name: "provider", Flag: "p", Description: "AI Provider", Default: "openai"},
		{Name: "language", Flag: "l", Description: "Language", Default: "pt_BR"},
	}
}

func (g GenerateCommit) Execute(input *cli.CommandInput) (*cli.ExecutionResult, error) {
	executionResult := cli.NewExecutionResult()
	executionResult.Message = fmt.Sprintf("My message here. {provider: %s, language: %s}",
		input.Options["provider"].Value, input.Options["language"].Value)
	return executionResult, nil
}
