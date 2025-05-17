package vo

import (
	"strings"
)

type ColoredText struct {
	text string
}

func NewColoredText(text string) *ColoredText {
	return &ColoredText{text: text}
}

func NewColoredMultilineText(lines []string) *ColoredText {
	joined := strings.Join(lines, "\n")
	return &ColoredText{text: joined}
}

func (c ColoredText) Render() string {
	replacer := strings.NewReplacer(
		"<info>", "\033[32m",
		"</info>", "\033[0m",
		"<error>", "\033[31m",
		"</error>", "\033[0m",
		"<comment>", "\033[33m",
		"</comment>", "\033[0m",
	)
	return replacer.Replace(c.text)
}
