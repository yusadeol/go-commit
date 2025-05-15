package main

import (
	"encoding/json"
	"fmt"
	"github.com/yusadeol/go-commit/internal/interface/cli"
	"github.com/yusadeol/go-commit/internal/interface/cli/command"
	"log"
	"os"
	"path/filepath"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		executionResult := cli.NewExecutionResult()
		executionResult.ExitCode = cli.ExitInvalidUsage
		executionResult.Message = "No command provided."
		printAndExit(executionResult)
	}
	configuration, err := loadConfiguration()
	if err != nil {
		executionResult := cli.NewExecutionResult()
		executionResult.ExitCode = cli.ExitError
		executionResult.Message = err.Error()
		printAndExit(executionResult)
	}
	commandsToRegister := []cli.Command{
		command.NewGenerateCommit(configuration),
	}
	commandDispatcher := cli.NewCommandDispatcher()
	for _, commandToRegister := range commandsToRegister {
		commandDispatcher.Register(commandToRegister)
	}
	executionResult, err := commandDispatcher.Dispatch(args[0], args[1:])
	if err != nil {
		executionResult = cli.NewExecutionResult()
		executionResult.ExitCode = cli.ExitError
		executionResult.Message = err.Error()
	}
	printAndExit(executionResult)
}

func loadConfiguration() (*command.Configuration, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	configPath := filepath.Join(homeDir, ".config", "commit.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	var config command.Configuration
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func printAndExit(executionResult *cli.ExecutionResult) {
	outputChannel := os.Stdout
	if executionResult.ExitCode != cli.ExitSuccess {
		outputChannel = os.Stderr
	}
	_, err := fmt.Fprintln(outputChannel, executionResult.Message)
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(int(executionResult.ExitCode))
}
