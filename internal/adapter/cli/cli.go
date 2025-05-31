package cli

import (
	"github.com/yusadeol/go-commit/internal/adapter/cli/dispatcher"
	"github.com/yusadeol/go-commit/internal/domain/vo"
)

type CLI struct {
	commandDispatcher *dispatcher.CommandDispatcher
}

func New(commandsToRegister []dispatcher.Command) *CLI {
	commandDispatcher := dispatcher.NewCommandDispatcher()
	for _, commandToRegister := range commandsToRegister {
		commandDispatcher.Register(commandToRegister)
	}
	return &CLI{commandDispatcher: commandDispatcher}
}

func (a CLI) Run(args []string) (*dispatcher.Result, error) {
	if len(args) == 0 {
		return &dispatcher.Result{
			ExitCode: vo.ExitCodeError,
			Message:  vo.NewMarkupText("<error>No command provided.</error>"),
		}, nil
	}
	return a.commandDispatcher.Dispatch(args[0], args[1:])
}
