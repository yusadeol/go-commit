package command

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"

	"github.com/yusadeol/go-commit/internal/Domain/vo"
	"github.com/yusadeol/go-commit/internal/application/usecase"
	"github.com/yusadeol/go-commit/internal/infrastructure/service/ai"
	"github.com/yusadeol/go-commit/internal/interface/cli"
)

type Generate struct {
	configuration            *vo.Configuration
	aiDefaultProviderFactory ai.ProviderFactory
}

func NewGenerate(configuration *vo.Configuration, aiDefaultProviderFactory ai.ProviderFactory) *Generate {
	return &Generate{configuration: configuration, aiDefaultProviderFactory: aiDefaultProviderFactory}
}

func (g *Generate) GetName() string {
	return "generate"
}

func (g *Generate) GetArguments() []cli.Argument {
	return []cli.Argument{
		{Name: "diff", Description: "Git diff", Required: false},
	}
}

func (g *Generate) GetOptions() []cli.Option {
	return []cli.Option{
		{Name: "provider", Flag: "p", Description: "AI Provider", Default: "openai"},
		{Name: "language", Flag: "l", Description: "Language", Default: "en_US"},
		{Name: "commit", Flag: "c", Description: "Commit", Default: "true"},
	}
}

func (g *Generate) Execute(input *cli.CommandInput) (*cli.Result, error) {
	Result := cli.NewResult()
	configurationAIProvider, configurationAIProviderExists := g.configuration.AIProviders[input.Options["provider"].Value]
	if !configurationAIProviderExists {
		return nil, fmt.Errorf("AI provider %q configuration not found", input.Options["provider"].Value)
	}
	configurationLanguage, configurationLanguageExists := g.configuration.Languages[input.Options["language"].Value]
	if !configurationLanguageExists {
		return nil, fmt.Errorf("language %q configuration not found", input.Options["language"].Value)
	}
	diff := input.Arguments["diff"].Value
	if diff == "" {
		var err error
		diff, err = g.GetGitDiff()
		if err != nil {
			return nil, err
		}
	}
	generate := usecase.NewGenerate()
	output, err := generate.Execute(&usecase.GenerateInput{
		AIDefaultProviderFactory: g.aiDefaultProviderFactory,
		AIProvider:               &configurationAIProvider,
		Language:                 &configurationLanguage,
		Diff:                     diff,
	})
	if err != nil {
		return nil, err
	}
	if input.Options["commit"].Value == "true" {
		err = g.CommitChanges(output.Commit)
		if err != nil {
			return nil, err
		}
	}
	message := []string{
		"<info>Commit generated and applied successfully!</info>",
		fmt.Sprintf("<comment>%s</comment>", output.Commit),
	}
	Result.Message = vo.NewColoredMultilineText(message)
	return Result, nil
}

func (g *Generate) GetGitDiff() (string, error) {
	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd := exec.Command("git", "diff", "--staged")
	cmd.Stdout = &out
	cmd.Stderr = &outErr
	err := cmd.Run()
	if err != nil {
		if outErr.Len() > 0 {
			return "", errors.New(outErr.String())
		}
		return "", err
	}
	diff := out.String()
	if diff == "" {
		return "", errors.New("no staged changes found")
	}
	return diff, nil
}

func (g *Generate) CommitChanges(commit string) error {
	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd := exec.Command("git", "commit", "-m", commit)
	cmd.Stdout = &out
	cmd.Stderr = &outErr
	err := cmd.Run()
	if err != nil {
		if outErr.Len() > 0 {
			return errors.New(outErr.String())
		}
		return err
	}
	return nil
}
