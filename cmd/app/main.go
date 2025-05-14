package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/yusadeol/go-commit/internal/interface/cli"
	"github.com/yusadeol/go-commit/internal/interface/cli/command"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	args := os.Args[1:]
	if len(args) == 0 {
		result := cli.NewExecutionResult()
		result.ExitCode = cli.ExitInvalidUsage
		result.Message = "No command provided."
		printAndExit(result)
	}
	commandsToRegister := []cli.Command{
		command.NewGenerateCommit(),
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
