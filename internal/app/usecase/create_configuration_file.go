package usecase

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/yusadeol/go-commit/internal/domain/vo"
)

var ErrConfigurationAlreadyExists = errors.New("configuration already exists")

var defaultConfiguration = vo.Configuration{
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

type CreateConfigurationFile struct{}

func NewCreateConfigurationFile() *CreateConfigurationFile {
	return &CreateConfigurationFile{}
}

func (c *CreateConfigurationFile) Execute(input *CreateConfigurationFileInput) error {
	_, err := os.Stat(input.ConfigurationDirPath)
	if err != nil {
		return err
	}
	configurationMarshal, err := json.MarshalIndent(defaultConfiguration, "", "    ")
	if err != nil {
		return err
	}
	configurationFilePath := filepath.Join(input.ConfigurationDirPath, "commit.json")
	_, err = os.Stat(configurationFilePath)
	if err == nil {
		return ErrConfigurationAlreadyExists
	}
	err = os.WriteFile(configurationFilePath, configurationMarshal, 0644)
	if err != nil {
		return err
	}
	return nil
}

type CreateConfigurationFileInput struct {
	ConfigurationDirPath string
}
