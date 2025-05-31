package command

import (
	"fmt"

	"github.com/yusadeol/go-commit/internal/adapter/cli/dispatcher"

	"github.com/yusadeol/go-commit/internal/domain/vo"
)

type Version struct {
	version string
}

func NewVersion(version string) *Version {
	return &Version{version: version}
}

func (g *Version) GetName() string {
	return "version"
}

func (g *Version) GetArguments() []dispatcher.Argument {
	return []dispatcher.Argument{}
}

func (g *Version) GetOptions() []dispatcher.Option {
	return []dispatcher.Option{}
}

func (g *Version) Execute(input *dispatcher.CommandInput) (*dispatcher.Result, error) {
	result := dispatcher.NewResult()
	result.Message = vo.NewMarkupText(fmt.Sprintf("<success>%s</success>", g.version))
	return result, nil
}
