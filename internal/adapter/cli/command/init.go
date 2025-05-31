package command

import (
	"errors"

	"github.com/yusadeol/go-commit/internal/adapter/cli/dispatcher"

	"github.com/yusadeol/go-commit/internal/app/usecase"
	"github.com/yusadeol/go-commit/internal/domain/vo"
)

type Init struct {
	configurationDirPath string
}

func NewInit(configurationDirPath string) *Init {
	return &Init{configurationDirPath: configurationDirPath}
}

func (g *Init) GetName() string {
	return "init"
}

func (g *Init) GetArguments() []dispatcher.Argument {
	return []dispatcher.Argument{}
}

func (g *Init) GetOptions() []dispatcher.Option {
	return []dispatcher.Option{}
}

func (g *Init) Execute(input *dispatcher.CommandInput) (*dispatcher.Result, error) {
	result := dispatcher.NewResult()
	createConfigurationFile := usecase.NewCreateConfigurationFile()
	err := createConfigurationFile.Execute(&usecase.CreateConfigurationFileInput{
		ConfigurationDirPath: g.configurationDirPath,
	})
	if errors.Is(err, usecase.ErrConfigurationAlreadyExists) {
		result.ExitCode = vo.ExitCodeError
		result.Message = vo.NewMarkupText("<info>configuration file already exists</info>")
		return result, nil
	}
	if err != nil {
		return nil, err
	}
	result.Message = vo.NewMarkupText("<success>configuration file created successfully</success>")
	return result, nil
}
