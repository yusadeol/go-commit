package usecase

import (
	"fmt"
	"github.com/yusadeol/go-commit/internal/infrastructure/service/ai"
)

type GenerateCommit struct{}

func NewGenerateCommit() *GenerateCommit {
	return &GenerateCommit{}
}

func (g *GenerateCommit) Execute(input *GenerateCommitInput) (*GenerateCommitOutput, error) {
	instructions := fmt.Sprintf(`
        Write a commit message for this diff following Conventional Commits specification.
		Do NOT use scopes.
		EACH line must not exceed 72 characters.
		Write the commit message in %s language without any accents.
		ONLY return the commit message, without any additional text or explanation.
		If there are multiple modifications in different contexts, write the body using a list format.
		Otherwise, use a regular paragraph format that ends with a period.
		If the body is a list, DO NOT add a period at the end of each list item, as in the following example:
		feat: add a new feature

		- Add a new feature
		- Fix a bug
	`, input.Language)
	output, err := input.Provider.Ask(ai.ProviderInput{Model: "gpt-4.1", Instructions: instructions, Input: input.Diff})
	if err != nil {
		return nil, err
	}
	return NewGenerateCommitOutput(output.Text), nil
}

type GenerateCommitInput struct {
	Provider ai.Provider
	Language string
	Diff     string
}

func NewGenerateCommitInput(provider ai.Provider, language string, diff string) *GenerateCommitInput {
	return &GenerateCommitInput{Provider: provider, Language: language, Diff: diff}
}

type GenerateCommitOutput struct {
	Commit string
}

func NewGenerateCommitOutput(commit string) *GenerateCommitOutput {
	return &GenerateCommitOutput{Commit: commit}
}
