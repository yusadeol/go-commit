package main

import (
	"encoding/json"
	"fmt"
	"github.com/yusadeol/go-commit/internal/Domain/vo"
	"github.com/yusadeol/go-commit/internal/interface/cli"
	"github.com/yusadeol/go-commit/internal/interface/cli/command"
	"os"
	"path/filepath"
)

func main() {
	args := os.Args[1:]
	config, err := loadConfiguration()
	if err != nil {
		exitWithMessage(vo.ExitCodeError, err.Error())
	}
	commandsToRegister := []cli.Command{
		command.NewGenerate(config),
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
	configPath := filepath.Join(homeDir, ".config", "commit.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	var config vo.Configuration
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func exitWithMessage(exitCode vo.ExitCode, message string) {
	outputChannel := os.Stdout
	if exitCode != vo.ExitCodeSuccess {
		outputChannel = os.Stderr
	}
	fmt.Fprintln(outputChannel, message)
	os.Exit(int(exitCode))
}
