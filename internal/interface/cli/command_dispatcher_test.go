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
		{Name: "first", Description: "My first argument", Required: true},
	}
}

func (m *mockCommand) GetOptions() []Option {
	return []Option{
		{Name: "first", Flag: "f", Description: "My first option", Default: "default-value"},
	}
}

func (m *mockCommand) Execute(input *CommandInput) (*Result, error) {
	m.executed = true
	m.input = input
	return NewResult(), nil
}

func (m *mockCommand) GetName() string {
	return "mock"
}

func TestCommandDispatcher(t *testing.T) {
	t.Run("dispatches a command with required argument and option", func(t *testing.T) {
		args := []string{"argumentValue", "--first", "optionValue"}
		command := newMockCommand()
		dispatcher := NewCommandDispatcher()
		dispatcher.Register(command)
		output, err := dispatcher.Dispatch("mock", args)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if output.ExitCode != ExitCodeSuccess {
			t.Fatalf("expected ExitCodeSuccess, got: %v", output.ExitCode)
		}
		if !command.executed {
			t.Fatal("expected command to be executed")
		}
		argumentFirst := command.input.Arguments["first"]
		if argumentFirst.Value != "argumentValue" {
			t.Errorf("expected argument 'first' to be 'argumentValue', got: %s", argumentFirst.Value)
		}
		optionFirst := command.input.Options["first"]
		if optionFirst.Value != "optionValue" {
			t.Errorf("expected option 'first' to be 'optionValue', got: %s", optionFirst.Value)
		}
	})

	t.Run("returns error when required argument is missing", func(t *testing.T) {
		args := []string{"--first", "optionValue"}
		command := newMockCommand()
		dispatcher := NewCommandDispatcher()
		dispatcher.Register(command)
		output, err := dispatcher.Dispatch("mock", args)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if output.ExitCode != ExitCodeInvalidUsage {
			t.Fatalf("expected ExitCodeInvalidUsage, got: %v", output.ExitCode)
		}
		expectedMsg := "Missing required argument(s): [first]"
		if output.Message.Render() != expectedMsg {
			t.Errorf("expected message %q, got: %q", expectedMsg, output.Message.Render())
		}
	})

	t.Run("returns error when command is not found", func(t *testing.T) {
		dispatcher := NewCommandDispatcher()
		output, err := dispatcher.Dispatch("notfound", []string{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if output.ExitCode != ExitCodeCommandNotFound {
			t.Fatalf("expected ExitCodeCommandNotFound, got: %v", output.ExitCode)
		}
	})

	t.Run("uses default option value when option is omitted", func(t *testing.T) {
		args := []string{"argumentValue"}
		command := newMockCommand()
		dispatcher := NewCommandDispatcher()
		dispatcher.Register(command)
		output, err := dispatcher.Dispatch("mock", args)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if output.ExitCode != ExitCodeSuccess {
			t.Fatalf("expected ExitCodeSuccess, got: %v", output.ExitCode)
		}
		opt := command.input.Options["first"]
		if opt.Value != "default-value" {
			t.Errorf("expected option to use default value 'default-value', got: %s", opt.Value)
		}
	})

	t.Run("returns error on unknown flag", func(t *testing.T) {
		args := []string{"argumentValue", "--unknown", "oops"}
		command := newMockCommand()
		dispatcher := NewCommandDispatcher()
		dispatcher.Register(command)
		output, err := dispatcher.Dispatch("mock", args)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if output.ExitCode != ExitCodeInvalidUsage {
			t.Fatalf("expected ExitCodeInvalidUsage, got: %v", output.ExitCode)
		}
	})
}
