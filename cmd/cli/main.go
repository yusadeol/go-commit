package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/yusadeol/go-commit/internal/adapter/cli/dispatcher"

	"github.com/yusadeol/go-commit/internal/adapter/cli"
	"github.com/yusadeol/go-commit/internal/adapter/cli/command"

	"github.com/yusadeol/go-commit/internal/domain/vo"
	"github.com/yusadeol/go-commit/internal/infra/service/ai"
)

func main() {
	args := os.Args[1:]
	configurationDirPath, err := getConfigurationFilePath()
	if err != nil {
		exitWithMessage(
			vo.ExitCodeError,
			vo.NewMarkupText(fmt.Sprintf("<error>%s</error>", err.Error())),
		)
	}
	configuration, err := loadConfiguration(configurationDirPath)
	if err != nil {
		exitWithMessage(
			vo.ExitCodeError,
			vo.NewMarkupText(fmt.Sprintf("<error>%s</error>", err.Error())),
		)
	}
	commandsToRegister := []dispatcher.Command{
		command.NewVersion("v1.0.1"),
		command.NewInit(configurationDirPath),
		command.NewGenerate(configuration, ai.NewDefaultProviderFactory()),
	}
	app := cli.New(commandsToRegister)
	output, err := app.Run(args)
	if err != nil {
		exitWithMessage(
			vo.ExitCodeError,
			vo.NewMarkupText(fmt.Sprintf("<error>%s</error>", err.Error())),
		)
	}
	exitWithMessage(output.ExitCode, output.Message)
}

func getConfigurationFilePath() (string, error) {
	homeDirPath, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDirPath, ".config"), nil
}

func exitWithMessage(exitCode vo.ExitCode, message *vo.MarkupText) {
	outputChannel := os.Stdout
	if exitCode != vo.ExitCodeSuccess {
		outputChannel = os.Stderr
	}
	_, _ = fmt.Fprintln(outputChannel, message.ToANSI())
	os.Exit(int(exitCode))
}

func loadConfiguration(configurationDirPath string) (*vo.Configuration, error) {
	configurationFilePath := filepath.Join(configurationDirPath, "commit.json")
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
