package cli

import (
	"testing"
)

type mockCommand struct {
	executed bool
	input    *CommandInput
}

func newMockCommand() *mockCommand {
	return &mockCommand{}
}

func (m *mockCommand) GetArguments() []Argument {
	return []Argument{
		{Name: "first", Description: "My first argument", Required: false},
	}
}

func (m *mockCommand) GetOptions() []Option {
	return []Option{
		{Name: "first", Flag: "f", Description: "My first option", Default: "default-value"},
	}
}

func (m *mockCommand) Execute(input *CommandInput) (*ExecutionResult, error) {
	m.executed = true
	m.input = input
	return NewExecutionResult(), nil
}

func (m *mockCommand) GetName() string {
	return "mock"
}

func TestNewCommandDispatcher(t *testing.T) {
	t.Run("should be able to dispatch a command", func(t *testing.T) {
		args := []string{
			"argumentValue",
			"--first",
			"optionValue",
		}
		command := newMockCommand()
		dispatcher := NewCommandDispatcher()
		dispatcher.Register(command)
		output, err := dispatcher.Dispatch("mock", args)
		if err != nil {
			t.Fatalf("Dispatch returned an error: %v", err)
		}
		if output.ExitCode != ExitSuccess {
			t.Fatalf("Expected ExitSuccess, got: %v", output.ExitCode)
		}
		if !command.executed {
			t.Fatal("Expected command to be executed")
		}
		argument, argumentExists := command.input.Arguments["first"]
		if !argumentExists || argument.Value != "argumentValue" {
			t.Errorf("Expected argument 'first' to be 'argumentValue', got: %+v", argument)
		}
		option, optionExists := command.input.Options["first"]
		if !optionExists || option.Value != "optionValue" {
			t.Errorf("Expected option 'first' to be 'optionValue', got: %+v", option)
		}
	})

	t.Run("should return command not found", func(t *testing.T) {
		dispatcher := NewCommandDispatcher()
		output, err := dispatcher.Dispatch("notfound", []string{})
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if output.ExitCode != ExitCommandNotFound {
			t.Fatalf("Expected ExitCommandNotFound, got: %v", output.ExitCode)
		}
	})
}
