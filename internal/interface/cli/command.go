package cli

type Command interface {
	GetName() string
	GetArguments() []Argument
	GetOptions() []Option
	Execute(input *CommandInput) (*Result, error)
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
