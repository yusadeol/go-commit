package main

import (
	"encoding/json"
	"fmt"
	"github.com/yusadeol/go-commit/internal/Domain/vo"
	"github.com/yusadeol/go-commit/internal/interface/cli"
	"github.com/yusadeol/go-commit/internal/interface/cli/command"
	"log"
	"os"
	"path/filepath"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		Result := cli.NewResult()
		Result.ExitCode = cli.ExitCodeInvalidUsage
		Result.Message = "No command provided."
		printAndExitCode(Result)
	}
	configuration, err := loadConfiguration()
	if err != nil {
		Result := cli.NewResult()
		Result.ExitCode = cli.ExitCodeError
		Result.Message = err.Error()
		printAndExitCode(Result)
	}
	commandsToRegister := []cli.Command{
		command.NewGenerateCommit(configuration),
	}
	commandDispatcher := cli.NewCommandDispatcher()
	for _, commandToRegister := range commandsToRegister {
		commandDispatcher.Register(commandToRegister)
	}
	Result, err := commandDispatcher.Dispatch(args[0], args[1:])
	if err != nil {
		Result = cli.NewResult()
		Result.ExitCode = cli.ExitCodeError
		Result.Message = err.Error()
	}
	printAndExitCode(Result)
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

func printAndExitCode(Result *cli.Result) {
	outputChannel := os.Stdout
	if Result.ExitCode != cli.ExitCodeSuccess {
		outputChannel = os.Stderr
	}
	_, err := fmt.Fprintln(outputChannel, Result.Message)
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(int(Result.ExitCode))
}
