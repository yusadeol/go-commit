package cli

import "strings"

type CommandDispatcher struct {
	commands map[string]Command
}

type Command interface {
	GetName() string
	GetArguments() []*Argument
	GetOptions() []*Option
	Execute(input *CommandInput) (*ExecutionResult, error)
}

type Argument struct {
	Name        string
	Description string
	Value       string
	Required    bool
}

type Option struct {
	Name        string
	Flag        string
	Description string
	Value       string
	Default     string
}

type CommandInput struct {
	Arguments map[string]*Argument
	Options   map[string]*Option
}

func NewCommandInput(arguments map[string]*Argument, options map[string]*Option) *CommandInput {
	return &CommandInput{Arguments: arguments, Options: options}
}

type ExecutionResult struct {
	ExitCode ExitCode
	Message  string
}

type ExitCode int

const (
	ExitSuccess           ExitCode = 0
	ExitError             ExitCode = 1
	ExitInvalidUsage      ExitCode = 2
	ExitCommandNotFound   ExitCode = 127
	ExitPermissionDenied  ExitCode = 126
	ExitInterruptedByUser ExitCode = 130
	ExitOutOfMemory       ExitCode = 137
)

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
	defaultExecutionResult := NewExecutionResult()
	command, exists := c.commands[calledCommandName]
	if !exists {
		defaultExecutionResult.ExitCode = ExitCommandNotFound
		return defaultExecutionResult, nil
	}
	commandInput, err := c.parseCommandInput(command.GetArguments(), command.GetOptions(), args)
	if err != nil {
		return nil, err
	}
	return command.Execute(commandInput)
}

func (c *CommandDispatcher) parseCommandInput(commandArguments []*Argument, commandOptions []*Option, args []string) (*CommandInput, error) {
	c.parseArguments(commandArguments, args)
	c.parseOptions(commandOptions, args)
	parsedArguments := make(map[string]*Argument)
	for _, argument := range commandArguments {
		parsedArguments[argument.Name] = argument
	}
	parsedOptions := make(map[string]*Option)
	for _, option := range commandOptions {
		parsedOptions[option.Name] = option
	}
	return NewCommandInput(parsedArguments, parsedOptions), nil
}

func (c *CommandDispatcher) parseArguments(commandArguments []*Argument, args []string) {
	argumentsFromArgs := c.getArgumentsFromArgs(args)
	for index, argumentArg := range argumentsFromArgs {
		if index >= len(commandArguments) {
			break
		}
		argument := commandArguments[index]
		argument.Value = argumentArg
	}
}

func (c *CommandDispatcher) getArgumentsFromArgs(args []string) []string {
	var argumentsFromArgs []string
	var previousArgIsOption bool
	for _, arg := range args {
		if previousArgIsOption {
			previousArgIsOption = false
			continue
		}
		if strings.HasPrefix(arg, "-") || strings.HasPrefix(arg, "--") {
			previousArgIsOption = true
			continue
		}
		argumentsFromArgs = append(argumentsFromArgs, arg)
	}
	return argumentsFromArgs
}

func (c *CommandDispatcher) parseOptions(commandOptions []*Option, args []string) {
	optionsFromArgs := c.getOptionsFromArgs(args)
	for index, optionArg := range optionsFromArgs {
		option := commandOptions[index]
		parts := strings.SplitN(optionArg, "=", 2)
		if len(parts) != 2 {
			continue
		}
		optionName := strings.TrimLeft(parts[0], "-")
		optionValue := parts[1]
		if option.Name == optionName || option.Flag == optionName {
			option.Value = optionValue
		}
	}
}

func (c *CommandDispatcher) getOptionsFromArgs(args []string) []string {
	var optionsFromArgs []string
	var previousArg string
	for _, arg := range args {
		if previousArg != "" {
			optionsFromArgs = append(optionsFromArgs, previousArg+"="+arg)
			previousArg = ""
			continue
		}
		if strings.HasPrefix(arg, "-") || strings.HasPrefix(arg, "--") {
			if strings.Contains(arg, "=") {
				optionsFromArgs = append(optionsFromArgs, arg)
				continue
			}
			previousArg = arg
		}
	}
	return optionsFromArgs
}
