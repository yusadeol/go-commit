package cli

import "github.com/yusadeol/go-commit/internal/domain/vo"

type CLI struct {
	commandDispatcher *CommandDispatcher
}

func New(commandsToRegister []Command) *CLI {
	commandDispatcher := NewCommandDispatcher()
	for _, commandToRegister := range commandsToRegister {
		commandDispatcher.Register(commandToRegister)
	}
	return &CLI{commandDispatcher: commandDispatcher}
}

func (a CLI) Run(args []string) (*Result, error) {
	if len(args) == 0 {
		return &Result{
			ExitCode: vo.ExitCodeError,
			Message:  vo.NewMarkupText("<error>No command provided.</error>"),
		}, nil
	}
	return a.commandDispatcher.Dispatch(args[0], args[1:])
}
