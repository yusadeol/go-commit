package cli

import "testing"

type mockCommand struct{}

func newMockCommand() *mockCommand {
	return &mockCommand{}
}

func (m *mockCommand) Execute() (*ExecutionResult, error) {
	return NewExecutionResult(), nil
}

func (m *mockCommand) GetName() string {
	return "mock"
}

func TestNewCommandDispatcher(t *testing.T) {
	t.Run("should be able to dispatch a command", func(t *testing.T) {
		args := []string{
			"--provider",
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
