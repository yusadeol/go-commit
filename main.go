package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/yusadeol/go-commit/internal/adapter/cli"
	"github.com/yusadeol/go-commit/internal/adapter/cli/command"

	"github.com/yusadeol/go-commit/internal/domain/vo"
	"github.com/yusadeol/go-commit/internal/infra/service/ai"
)

func main() {
	args := os.Args[1:]
	homeDirPath, err := os.UserHomeDir()
	if err != nil {
		exitWithMessage(
			vo.ExitCodeError,
			vo.NewMarkupText(fmt.Sprintf("<error>%s</error>", err.Error())),
		)
	}
	configurationDirPath := filepath.Join(homeDirPath, ".config")
	configuration, err := loadConfiguration(configurationDirPath)
	if err != nil {
		exitWithMessage(
			vo.ExitCodeError,
			vo.NewMarkupText(fmt.Sprintf("<error>%s</error>", err.Error())),
		)
	}
	commandsToRegister := []cli.Command{
		command.NewInit(configurationDirPath),
		command.NewGenerate(configuration, ai.NewDefaultProviderFactory()),
	}
	app := cli.NewApplication(commandsToRegister)
	output, err := app.Run(args)
	if err != nil {
		exitWithMessage(
			vo.ExitCodeError,
			vo.NewMarkupText(fmt.Sprintf("<error>%s</error>", err.Error())),
		)
	}
	exitWithMessage(output.ExitCode, output.Message)
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
