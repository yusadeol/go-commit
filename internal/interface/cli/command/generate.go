package command

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/yusadeol/go-commit/internal/Domain/vo"
	"github.com/yusadeol/go-commit/internal/application/usecase"
	"github.com/yusadeol/go-commit/internal/infrastructure/service/ai"
	"github.com/yusadeol/go-commit/internal/interface/cli"
	"os/exec"
)

type Generate struct {
	configuration *vo.Configuration
}

func NewGenerate(configuration *vo.Configuration) *Generate {
	return &Generate{configuration: configuration}
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
	}
}

func (g *Generate) Execute(input *cli.CommandInput) (*cli.Result, error) {
	Result := cli.NewResult()
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
	generateInput := usecase.NewGenerateInput(
		configurationAiProvider.DefaultModel,
		aiProvider,
		input.Options["language"].Value,
		diff,
	)
	generate := usecase.NewGenerate()
	output, err := generate.Execute(generateInput)
	if err != nil {
		return nil, err
	}
	err = g.CommitChanges(output.Commit)
	if err != nil {
		return nil, err
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

func (g *Generate) CommitChanges(commit string) error {
	var out bytes.Buffer
	cmd := exec.Command("git", "commit", "-m", commit)
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
