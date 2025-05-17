package cli

import "github.com/yusadeol/go-commit/internal/Domain/vo"

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

func (a Application) Run(args []string) (*ApplicationOutput, error) {
	if len(args) == 0 {
		return &ApplicationOutput{
			ExitCode: vo.ExitCodeError,
			Message:  vo.NewColoredText("<error>No command provided.</error>").Render(),
		}, nil
	}
	result, err := a.commandDispatcher.Dispatch(args[0], args[1:])
	if err != nil {
		return nil, err
	}
	return &ApplicationOutput{
		ExitCode: result.ExitCode,
		Message:  result.Message.Render(),
	}, nil
}

type ApplicationOutput struct {
	ExitCode vo.ExitCode
	Message  string
}
