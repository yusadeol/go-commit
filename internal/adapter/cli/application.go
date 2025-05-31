package cli

import "github.com/yusadeol/go-commit/internal/domain/vo"

type Application struct {
	commandDispatcher *CommandDispatcher
}

func NewApplication(commandsToRegister []Command) *Application {
	commandDispatcher := NewCommandDispatcher()
	for _, commandToRegister := range commandsToRegister {
		commandDispatcher.Register(commandToRegister)
	}
	return &Application{commandDispatcher: commandDispatcher}
}

func (a Application) Run(args []string) (*Result, error) {
	if len(args) == 0 {
		return &Result{
			ExitCode: vo.ExitCodeError,
			Message:  vo.NewMarkupText("<error>No command provided.</error>"),
		}, nil
	}
	return a.commandDispatcher.Dispatch(args[0], args[1:])
}
