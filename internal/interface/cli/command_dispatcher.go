package cli

type CommandDispatcher struct {
	commands map[string]Command
}

type Command interface {
	GetName() string
	Execute() (*ExecutionResult, error)
}

type ExitCode int

const (
	ExitSuccess      ExitCode = 0
	ExitError        ExitCode = 1
	ExitInvalidUsage ExitCode = 2
	ExitCommandNF    ExitCode = 127
	ExitPermission   ExitCode = 126
	ExitInterrupted  ExitCode = 130
	ExitOutOfMemory  ExitCode = 137
)

type ExecutionResult struct {
	ExitCode ExitCode
	Message  string
}

func NewExecutionResult() *ExecutionResult {
	return &ExecutionResult{ExitCode: ExitSuccess}
}

func NewCommandDispatcher() *CommandDispatcher {
	return &CommandDispatcher{commands: make(map[string]Command)}
}

func (c *CommandDispatcher) Register(command Command) {
	c.commands[command.GetName()] = command
}

func (c *CommandDispatcher) Dispatch(calledCommandName string, args []string) (*ExecutionResult, error) {
	return NewExecutionResult(), nil
}
