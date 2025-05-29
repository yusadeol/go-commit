package main

import (
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
	homeDirPath, err := os.UserHomeDir()
	if err != nil {
		exitWithMessage(
			vo.ExitCodeError,
			vo.NewMarkupText(fmt.Sprintf("<error>%s</error>", err.Error())),
		)
	}
	configurationDirPath := filepath.Join(homeDirPath, ".config")
	commandsToRegister := []cli.Command{
		command.NewInit(configurationDirPath),
		command.NewGenerate(ai.NewDefaultProviderFactory(), configurationDirPath),
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
	fmt.Fprintln(outputChannel, message.ToANSI())
	os.Exit(int(exitCode))
}
