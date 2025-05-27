package cli

import (
	"testing"

	"github.com/yusadeol/go-commit/internal/Domain/vo"
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
		{
			Name:          "first",
			Flag:          "f",
			Description:   "My first option",
			AllowedValues: []string{"option-value", "default-value"},
			Default:       "default-value",
		},
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
		args := []string{"argument-value", "--first", "option-value"}
		command := newMockCommand()
		dispatcher := NewCommandDispatcher()
		dispatcher.Register(command)
		output, err := dispatcher.Dispatch("mock", args)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if output.ExitCode != vo.ExitCodeSuccess {
			t.Fatalf("expected ExitCodeSuccess, got: %v", output.ExitCode)
		}
		if !command.executed {
			t.Fatal("expected command to be executed")
		}
		argumentFirst := command.input.Arguments["first"]
		if argumentFirst.Value != "argument-value" {
			t.Errorf("expected argument 'first' to be 'argument-value', got: %s", argumentFirst.Value)
		}
		optionFirst := command.input.Options["first"]
		if optionFirst.Value != "option-value" {
			t.Errorf("expected option 'first' to be 'option-value', got: %s", optionFirst.Value)
		}
	})

	t.Run("returns error when required argument is missing", func(t *testing.T) {
		args := []string{"--first", "option-value"}
		command := newMockCommand()
		dispatcher := NewCommandDispatcher()
		dispatcher.Register(command)
		output, err := dispatcher.Dispatch("mock", args)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if output.ExitCode != vo.ExitCodeInvalidUsage {
			t.Fatalf("expected ExitCodeInvalidUsage, got: %v", output.ExitCode)
		}
		expectedMessage := "missing required argument: first"
		gotMessage := output.Message.Plain()
		if gotMessage != expectedMessage {
			t.Errorf("expected message %q, got: %q", expectedMessage, gotMessage)
		}
	})

	t.Run("returns error when command is not found", func(t *testing.T) {
		dispatcher := NewCommandDispatcher()
		output, err := dispatcher.Dispatch("notfound", []string{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if output.ExitCode != vo.ExitCodeCommandNotFound {
			t.Fatalf("expected ExitCodeCommandNotFound, got: %v", output.ExitCode)
		}
	})

	t.Run("uses default option value when option is omitted", func(t *testing.T) {
		args := []string{"argument-value"}
		command := newMockCommand()
		dispatcher := NewCommandDispatcher()
		dispatcher.Register(command)
		output, err := dispatcher.Dispatch("mock", args)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if output.ExitCode != vo.ExitCodeSuccess {
			t.Fatalf("expected ExitCodeSuccess, got: %v", output.ExitCode)
		}
		opt := command.input.Options["first"]
		if opt.Value != "default-value" {
			t.Errorf("expected option to use default value 'default-value', got: %s", opt.Value)
		}
	})

	t.Run("returns error on unknown flag", func(t *testing.T) {
		args := []string{"argument-value", "--unknown", "oops"}
		command := newMockCommand()
		dispatcher := NewCommandDispatcher()
		dispatcher.Register(command)
		output, err := dispatcher.Dispatch("mock", args)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if output.ExitCode != vo.ExitCodeInvalidUsage {
			t.Fatalf("expected ExitCodeInvalidUsage, got: %v", output.ExitCode)
		}
	})

	t.Run("returns error on not allowed option value", func(t *testing.T) {
		args := []string{"--first", "not-allowed"}
		command := newMockCommand()
		dispatcher := NewCommandDispatcher()
		dispatcher.Register(command)
		output, err := dispatcher.Dispatch("mock", args)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if output.ExitCode != vo.ExitCodeInvalidUsage {
			t.Fatalf("expected ExitCodeInvalidUsage, got: %v", output.ExitCode)
		}
	})
}
