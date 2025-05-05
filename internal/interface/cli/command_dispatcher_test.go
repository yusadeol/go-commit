package cli

import (
	"testing"
)

type mockCommand struct{}

func newMockCommand() *mockCommand {
	return &mockCommand{}
}

func (m *mockCommand) GetArguments() []*Argument {
	return []*Argument{
		{Name: "first", Description: "My first argument", Required: false},
	}
}

func (m *mockCommand) GetOptions() []*Option {
	return []*Option{
		{Name: "first", Flag: "f", Description: "My first option", Default: "value"},
	}
}

func (m *mockCommand) Execute(input *CommandInput) (*ExecutionResult, error) {
	return NewExecutionResult(), nil
}

func (m *mockCommand) GetName() string {
	return "mock"
}

func TestNewCommandDispatcher(t *testing.T) {
	t.Run("should be able to dispatch a command", func(t *testing.T) {
		args := []string{
			"--first",
			"ysocode",
		}
		command := newMockCommand()
		commandDispatcher := NewCommandDispatcher()
		commandDispatcher.Register(command)
		output, err := commandDispatcher.Dispatch("mock", args)
		if err != nil {
			t.Fatal(err)
		}
		if output.ExitCode != ExitSuccess {
			t.Fatal("")
		}
	})
}
