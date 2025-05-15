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

func NewCommandDispatcher() *CommandDispatcher {
	return &CommandDispatcher{commands: make(map[string]Command)}
}

func (c *CommandDispatcher) Register(command Command) {
	c.commands[command.GetName()] = command
}

func (c *CommandDispatcher) Dispatch(calledCommandName string, args []string) (*Result, error) {
	defaultResult := NewResult()
	command, exists := c.commands[calledCommandName]
	if !exists {
		defaultResult.ExitCode = ExitCodeCommandNotFound
		return defaultResult, nil
	}
	commandInput, err := c.parseCommandInput(command.GetArguments(), c.standardizeOptions(command.GetOptions()), args)
	if err != nil {
		return nil, err
	}
	if len(c.missingRequiredArguments) > 0 {
		defaultResult.ExitCode = ExitCodeInvalidUsage
		defaultResult.Message = fmt.Sprintf("Missing required argument(s): %v", c.missingRequiredArguments)
		return defaultResult, nil
	}
	if len(c.unknownOptions) > 0 {
		defaultResult.ExitCode = ExitCodeInvalidUsage
		defaultResult.Message = fmt.Sprintf("Unknown option(s): %v", c.unknownOptions)
		return defaultResult, nil
	}
	return command.Execute(commandInput)
}

func (c *CommandDispatcher) standardizeOptions(options []Option) map[string]Option {
	standardizedOptions := make(map[string]Option, len(options))
	for _, option := range options {
		standardizedOptions[option.Name] = option
	}
	return standardizedOptions
}

func (c *CommandDispatcher) parseCommandInput(arguments []Argument, options map[string]Option, args []string) (*CommandInput, error) {
	argumentInputs, err := c.parseArguments(arguments, args)
	if err != nil {
		return nil, err
	}
	optionInputs, err := c.parseOptions(options, args)
	if err != nil {
		return nil, err
	}
	return NewCommandInput(argumentInputs, optionInputs), nil
}

func (c *CommandDispatcher) parseArguments(arguments []Argument, args []string) (map[string]ArgumentInput, error) {
	argumentsFromArgs := c.extractArgumentsFromArgs(args)
	argumentInputs := make(map[string]ArgumentInput, len(argumentsFromArgs))
	for index, argument := range arguments {
		if index < len(argumentsFromArgs) {
			argumentInputs[argument.Name] = ArgumentInput{Value: argumentsFromArgs[index], Meta: argument}
			continue
		}
		if argument.Required {
			c.missingRequiredArguments = append(c.missingRequiredArguments, argument.Name)
		}
	}
	return argumentInputs, nil
}

func (c *CommandDispatcher) extractArgumentsFromArgs(args []string) []string {
	var argumentsFromArgs []string
	var previousArgIsOption bool
	for _, arg := range args {
		if previousArgIsOption {
			previousArgIsOption = false
			continue
		}
		if strings.HasPrefix(arg, "-") {
			previousArgIsOption = true
			continue
		}
		argumentsFromArgs = append(argumentsFromArgs, arg)
	}
	return argumentsFromArgs
}

func (c *CommandDispatcher) parseOptions(options map[string]Option, args []string) (map[string]OptionInput, error) {
	optionsFromArgs := c.extractOptionsFromArgs(args)
	standardizedOptionsFromArgs := c.standardizeOptionsFromArgs(optionsFromArgs)
	optionInputs := map[string]OptionInput{}
	recognizedOptions := map[string]bool{}
	for _, option := range options {
		standardizedOptionValue, standardizedOptionExists := standardizedOptionsFromArgs[option.Name]
		matchedOptionIdentifier := option.Name
		if !standardizedOptionExists {
			standardizedOptionValue, standardizedOptionExists = standardizedOptionsFromArgs[option.Flag]
			matchedOptionIdentifier = option.Flag
			if !standardizedOptionExists {
				optionInputs[option.Name] = OptionInput{Meta: option}
				continue
			}
		}
		recognizedOptions[matchedOptionIdentifier] = true
		optionInputs[option.Name] = OptionInput{Value: standardizedOptionValue, Meta: option}
	}
	for optionName, optionInput := range optionInputs {
		if optionInput.Value == "" {
			optionInput.Value = optionInput.Meta.Default
		}
		optionInputs[optionName] = optionInput
	}
	for standardizedOptionName := range standardizedOptionsFromArgs {
		if !recognizedOptions[standardizedOptionName] {
			c.unknownOptions = append(c.unknownOptions, standardizedOptionName)
		}
	}
	return optionInputs, nil
}

func (c *CommandDispatcher) extractOptionsFromArgs(args []string) []string {
	var optionsFromArgs []string
	var previousArg string
	for _, arg := range args {
		if previousArg != "" {
			optionsFromArgs = append(optionsFromArgs, previousArg+"="+arg)
			previousArg = ""
			continue
		}
		if !strings.HasPrefix(arg, "-") {
			continue
		}
		if !strings.Contains(arg, "=") {
			previousArg = arg
			continue
		}
		optionsFromArgs = append(optionsFromArgs, arg)
	}
	return optionsFromArgs
}

func (c *CommandDispatcher) standardizeOptionsFromArgs(optionsFromArgs []string) map[string]string {
	standardizedOptionsFromArgs := map[string]string{}
	for _, option := range optionsFromArgs {
		parts := strings.SplitN(option, "=", 2)
		if len(parts) != 2 {
			continue
		}
		optionName := strings.TrimLeft(parts[0], "-")
		optionValue := parts[1]
		standardizedOptionsFromArgs[optionName] = optionValue
	}
	return standardizedOptionsFromArgs
}
