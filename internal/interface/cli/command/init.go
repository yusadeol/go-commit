package command

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/yusadeol/go-commit/internal/Domain/vo"
	"github.com/yusadeol/go-commit/internal/interface/cli"
)

var initConfiguration = vo.Configuration{
	DefaultAIProvider: "openai",
	DefaultLanguage:   "en_US",
	AIProviders: map[string]vo.AIProvider{
		"openai": {
			ID:     "openai",
			APIKey: "",
			Models: []string{
				"gpt-4.1",
			},
			DefaultModel: "gpt-4.1",
		},
	},
	Languages: map[string]vo.Language{
		"en_US": {
			ID:          "en_US",
			DisplayName: "English (United States)",
		},
		"pt_BR": {
			ID:          "pt_BR",
			DisplayName: "Portuguese (Brazil)",
		},
		"es_ES": {
			ID:          "es_ES",
			DisplayName: "Spanish (Spain)",
		},
	},
}

type Init struct{}

func NewInit() *Init {
	return &Init{}
}

func (g *Init) GetName() string {
	return "init"
}

func (g *Init) GetArguments() []cli.Argument {
	return []cli.Argument{}
}

func (g *Init) GetOptions() []cli.Option {
	return []cli.Option{}
}

func (g *Init) Execute(input *cli.CommandInput) (*cli.Result, error) {
	result := cli.NewResult()
	homeDirPath, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	configurationDirPath := filepath.Join(homeDirPath, ".config")
	_, err = os.Stat(configurationDirPath)
	if err != nil {
		return nil, err
	}
	configurationMarshal, err := json.MarshalIndent(initConfiguration, "", "    ")
	if err != nil {
		return nil, err
	}
	configurationFilePath := filepath.Join(configurationDirPath, "commit.json")
	_, err = os.Stat(configurationFilePath)
	if err == nil {
		result.ExitCode = vo.ExitCodeError
		result.Message = vo.NewMarkupText("<info>configuration file already exists</info>")
		return result, nil
	}
	err = os.WriteFile(configurationFilePath, configurationMarshal, 0644)
	if err != nil {
		return nil, err
	}
	result.Message = vo.NewMarkupText("<success>configuration file created successfully</success>")
	return result, nil
}
