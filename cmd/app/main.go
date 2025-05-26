package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/yusadeol/go-commit/internal/Domain/vo"
	"github.com/yusadeol/go-commit/internal/infrastructure/service/ai"
	"github.com/yusadeol/go-commit/internal/interface/cli"
	"github.com/yusadeol/go-commit/internal/interface/cli/command"
)

func main() {
	args := os.Args[1:]
	configuration, err := loadConfiguration()
	if err != nil {
		exitWithMessage(vo.ExitCodeError, err.Error())
	}
	commandsToRegister := []cli.Command{
		command.NewGenerate(configuration, ai.NewDefaultProviderFactory()),
	}
	app := cli.NewApplication(commandsToRegister)
	output, err := app.Run(args)
	if err != nil {
		exitWithMessage(vo.ExitCodeError, err.Error())
	}
	exitWithMessage(output.ExitCode, output.Message)
}

func loadConfiguration() (*vo.Configuration, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	configurationPath := filepath.Join(homeDir, ".config", "commit.json")
	data, err := os.ReadFile(configurationPath)
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

func exitWithMessage(exitCode vo.ExitCode, message string) {
	outputChannel := os.Stdout
	if exitCode != vo.ExitCodeSuccess {
		outputChannel = os.Stderr
	}
	fmt.Fprintln(outputChannel, message)
	os.Exit(int(exitCode))
}
