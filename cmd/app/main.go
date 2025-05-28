package main

import (
	"fmt"
	"os"

	"github.com/yusadeol/go-commit/internal/Domain/vo"
	"github.com/yusadeol/go-commit/internal/infrastructure/service/ai"
	"github.com/yusadeol/go-commit/internal/interface/cli"
	"github.com/yusadeol/go-commit/internal/interface/cli/command"
)

func main() {
	args := os.Args[1:]
	commandsToRegister := []cli.Command{
		command.NewInit(),
		command.NewGenerate(ai.NewDefaultProviderFactory()),
	}
	app := cli.NewApplication(commandsToRegister)
	output, err := app.Run(args)
	if err != nil {
		exitWithMessage(vo.ExitCodeError, err.Error())
	}
	exitWithMessage(output.ExitCode, output.Message)
}

func exitWithMessage(exitCode vo.ExitCode, message string) {
	outputChannel := os.Stdout
	if exitCode != vo.ExitCodeSuccess {
		outputChannel = os.Stderr
	}
	fmt.Fprintln(outputChannel, message)
	os.Exit(int(exitCode))
}
