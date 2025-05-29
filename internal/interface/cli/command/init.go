package command

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/yusadeol/go-commit/internal/Domain/vo"
	"github.com/yusadeol/go-commit/internal/application/usecase"
	"github.com/yusadeol/go-commit/internal/interface/cli"
)

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
	createConfigurationFile := usecase.NewCreateConfigurationFile()
	err = createConfigurationFile.Execute(&usecase.CreateConfigurationFileInput{
		ConfigurationDirPath: configurationDirPath,
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
