package cli

import (
	"fmt"
	"strings"
)

type CommandDispatcher struct {
	commands                 map[string]Command
	missingRequiredArguments []string
	unknownOptions           []string
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
	if len(c.missingRequiredArguments) > 0 {
		defaultExecutionResult.ExitCode = ExitInvalidUsage
		defaultExecutionResult.Message = fmt.Sprintf("Missing required argument(s): %v", c.missingRequiredArguments)
		return defaultExecutionResult, nil
	}
	if len(c.unknownOptions) > 0 {
		defaultExecutionResult.ExitCode = ExitInvalidUsage
		defaultExecutionResult.Message = fmt.Sprintf("Unknown option(s): %v", c.unknownOptions)
		return defaultExecutionResult, nil
	}
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
	commandArgumentsInput, err := c.parseArguments(commandArguments, args)
	if err != nil {
		return nil, err
	}
	commandOptionsInput, err := c.parseOptions(commandOptions, args)
	if err != nil {
		return nil, err
	}
	return NewCommandInput(commandArgumentsInput, commandOptionsInput), nil
}

func (c *CommandDispatcher) parseArguments(commandArguments []Argument, args []string) (map[string]ArgumentInput, error) {
	argumentsFromArgs := c.getArgumentsFromArgs(args)
	commandArgumentsInput := map[string]ArgumentInput{}
	for index, argument := range commandArguments {
		if index < len(argumentsFromArgs) {
			value := argumentsFromArgs[index]
			commandArgumentsInput[argument.Name] = ArgumentInput{Value: value, Meta: argument}
			continue
		}
		if argument.Required {
			c.missingRequiredArguments = append(c.missingRequiredArguments, argument.Name)
		}
	}
	return commandArgumentsInput, nil
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

func (c *CommandDispatcher) parseOptions(commandOptions map[string]Option, args []string) (map[string]OptionInput, error) {
	optionsFromArgs := c.getOptionsFromArgs(args)
	standardizedOptionsFromArgs := c.standardizeOptionsFromArgs(optionsFromArgs)
	commandOptionsInput := map[string]OptionInput{}
	recognizedOptions := map[string]bool{}
	for _, option := range commandOptions {
		standardizedOptionValue, standardizedOptionExists := standardizedOptionsFromArgs[option.Name]
		matchedOptionIdentifier := option.Name
		if !standardizedOptionExists {
			standardizedOptionValue, standardizedOptionExists = standardizedOptionsFromArgs[option.Flag]
			matchedOptionIdentifier = option.Flag
			if !standardizedOptionExists {
				commandOptionsInput[option.Name] = OptionInput{Meta: option}
				continue
			}
		}
		recognizedOptions[matchedOptionIdentifier] = true
		commandOptionsInput[option.Name] = OptionInput{Value: standardizedOptionValue, Meta: option}
	}
	for commandOptionName, commandOptionInput := range commandOptionsInput {
		if commandOptionInput.Value == "" {
			commandOptionInput.Value = commandOptionInput.Meta.Default
		}
		commandOptionsInput[commandOptionName] = commandOptionInput
	}
	for standardizedOptionName := range standardizedOptionsFromArgs {
		if !recognizedOptions[standardizedOptionName] {
			c.unknownOptions = append(c.unknownOptions, standardizedOptionName)
		}
	}
	return commandOptionsInput, nil
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

func (c *CommandDispatcher) standardizeOptionsFromArgs(optionArguments []string) map[string]string {
	standardizedOptions := map[string]string{}
	for _, optionArgument := range optionArguments {
		parts := strings.SplitN(optionArgument, "=", 2)
		if len(parts) != 2 {
			continue
		}
		optionName := strings.TrimLeft(parts[0], "-")
		optionValue := parts[1]
		standardizedOptions[optionName] = optionValue
	}
	return standardizedOptions
}
