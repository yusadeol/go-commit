package cli

import (
	"strings"
)

type CommandDispatcher struct {
	commands map[string]Command
}

type Command interface {
	GetName() string
	GetArguments() []Argument
	GetOptions() []Option
	Execute(input *CommandInput) (*ExecutionResult, error)
}

type Argument struct {
	Name        string
	Description string
	Required    bool
}

type Option struct {
	Name        string
	Flag        string
	Description string
	Default     string
}

type CommandInput struct {
	Arguments map[string]ArgumentInput
	Options   map[string]OptionInput
}

type ArgumentInput struct {
	Value string
	Meta  Argument
}

type OptionInput struct {
	Value string
	Meta  Option
}

func NewCommandInput(arguments map[string]ArgumentInput, options map[string]OptionInput) *CommandInput {
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
	commandInput, err := c.parseCommandInput(command.GetArguments(), c.standardizeOptions(command.GetOptions()), args)
	if err != nil {
		return nil, err
	}
	return command.Execute(commandInput)
}

func (c *CommandDispatcher) standardizeOptions(commandOptions []Option) map[string]Option {
	standardizeOptions := map[string]Option{}
	for _, option := range commandOptions {
		standardizeOptions[option.Name] = option
	}
	return standardizeOptions
}

func (c *CommandDispatcher) parseCommandInput(commandArguments []Argument, commandOptions map[string]Option, args []string) (*CommandInput, error) {
	commandArgumentsInput := c.parseArguments(commandArguments, args)
	commandOptionsInput := c.parseOptions(commandOptions, args)
	return NewCommandInput(commandArgumentsInput, commandOptionsInput), nil
}

func (c *CommandDispatcher) parseArguments(commandArguments []Argument, args []string) map[string]ArgumentInput {
	argumentsFromArgs := c.getArgumentsFromArgs(args)
	commandArgumentsInput := map[string]ArgumentInput{}
	for index, argumentArg := range argumentsFromArgs {
		if index >= len(commandArguments) {
			break
		}
		argument := commandArguments[index]
		commandArgumentsInput[argument.Name] = ArgumentInput{Value: argumentArg, Meta: argument}
	}
	return commandArgumentsInput
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

func (c *CommandDispatcher) parseOptions(commandOptions map[string]Option, args []string) map[string]OptionInput {
	optionFlagsAndOptionNames := map[string]string{}
	for _, option := range commandOptions {
		optionFlagsAndOptionNames[option.Flag] = option.Name
	}
	commandOptionsInput := map[string]OptionInput{}
	optionsFromArgs := c.getOptionsFromArgs(args)
	for _, optionArg := range optionsFromArgs {
		parts := strings.SplitN(optionArg, "=", 2)
		if len(parts) != 2 {
			continue
		}
		optionName := strings.TrimLeft(parts[0], "-")
		if len(optionName) == 1 {
			mappedName, nameExists := optionFlagsAndOptionNames[optionName]
			if !nameExists {
				continue
			}
			optionName = mappedName
		}
		option, optionExists := commandOptions[optionName]
		if !optionExists {
			continue
		}
		commandOptionsInput[optionName] = OptionInput{
			Value: parts[1],
			Meta:  option,
		}
	}
	return commandOptionsInput
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
