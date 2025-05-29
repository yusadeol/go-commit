package command

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/yusadeol/go-commit/internal/Domain/vo"
	"github.com/yusadeol/go-commit/internal/application/usecase"
	"github.com/yusadeol/go-commit/internal/infrastructure/service/ai"
	"github.com/yusadeol/go-commit/internal/interface/cli"
)

type Generate struct {
	aiDefaultProviderFactory ai.ProviderFactory
	configurationDirPath     string
}

func NewGenerate(aiDefaultProviderFactory ai.ProviderFactory, configurationDirPath string) *Generate {
	return &Generate{aiDefaultProviderFactory: aiDefaultProviderFactory, configurationDirPath: configurationDirPath}
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
		{
			Name:          "provider",
			Flag:          "p",
			Description:   "AI Provider",
			AllowedValues: []string{"openai"},
			Default:       "openai",
		},
		{
			Name:          "language",
			Flag:          "l",
			Description:   "Language",
			AllowedValues: []string{"en_US", "pt_BR", "es_ES"},
			Default:       "en_US",
		},
		{
			Name:          "commit",
			Flag:          "c",
			Description:   "Commit",
			AllowedValues: []string{"true", "false"},
			Default:       "true",
		},
	}
}

func (g *Generate) Execute(input *cli.CommandInput) (*cli.Result, error) {
	result := cli.NewResult()
	configuration, err := g.loadConfiguration()
	if err != nil {
		return nil, err
	}
	configurationAIProvider, configurationAIProviderExists := configuration.AIProviders[input.Options["provider"].Value]
	if !configurationAIProviderExists {
		return nil, fmt.Errorf("AI provider %q configuration not found", input.Options["provider"].Value)
	}
	configurationLanguage, configurationLanguageExists := configuration.Languages[input.Options["language"].Value]
	if !configurationLanguageExists {
		return nil, fmt.Errorf("language %q configuration not found", input.Options["language"].Value)
	}
	diff := input.Arguments["diff"].Value
	if diff == "" {
		var err error
		diff, err = g.getGitDiff()
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
		err = g.commitChanges(output.Commit)
		if err != nil {
			return nil, err
		}
	}
	message := []string{
		"<info>Commit generated and applied successfully!</info>",
		fmt.Sprintf("<comment>%s</comment>", output.Commit),
	}
	result.Message = vo.NewColoredMultilineText(message)
	return result, nil
}

func (g *Generate) loadConfiguration() (*vo.Configuration, error) {
	configurationFilePath := filepath.Join(g.configurationDirPath, "commit.json")
	data, err := os.ReadFile(configurationFilePath)
	if err != nil {
		return nil, err
	}
	var configuration vo.Configuration
	err = json.Unmarshal(data, &configuration)
	if err != nil {
		return nil, err
	}
	return &configuration, nil
}

func (g *Generate) getGitDiff() (string, error) {
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

func (g *Generate) commitChanges(commit string) error {
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
