package usecase

import (
	"fmt"

	"github.com/yusadeol/go-commit/internal/Domain/vo"
	"github.com/yusadeol/go-commit/internal/infrastructure/service/ai"
)

type Generate struct{}

func NewGenerate() *Generate {
	return &Generate{}
}

func (g *Generate) Execute(input *GenerateInput) (*GenerateOutput, error) {
	aiProvider, err := input.ProviderFactory.Create(input.AIProvider.ID, input.AIProvider.APIKey)
	if err != nil {
		return nil, err
	}
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
	`, input.Language.DisplayName)
	output, err := aiProvider.Ask(&ai.ProviderInput{
		Model:        input.AIProvider.DefaultModel,
		Instructions: instructions,
		Input:        input.Diff,
	})
	if err != nil {
		return nil, err
	}
	return &GenerateOutput{Commit: output.Text}, nil
}

type GenerateInput struct {
	ProviderFactory *ai.ProviderFactory
	AIProvider      *vo.AIProvider
	Language        *vo.Language
	Diff            string
}

type GenerateOutput struct {
	Commit string
}
