package vo

import (
	"strings"
)

var markupToANSI = map[string]string{
	"<info>":     "\033[32m",
	"</info>":    "\033[0m",
	"<error>":    "\033[31m",
	"</error>":   "\033[0m",
	"<comment>":  "\033[33m",
	"</comment>": "\033[0m",
}

type MarkupText struct {
	text string
}

func NewMarkupText(text string) *MarkupText {
	return &MarkupText{text: text}
}

func NewColoredMultilineText(lines []string) *MarkupText {
	joined := strings.Join(lines, "\n")
	return &MarkupText{text: joined}
}

func (m MarkupText) ToANSI() string {
	out := m.text
	for tag, ansi := range markupToANSI {
		out = strings.ReplaceAll(out, tag, ansi)
	}
	return out
}

func (m MarkupText) Plain() string {
	out := m.text
	for tag := range markupToANSI {
		out = strings.ReplaceAll(out, tag, "")
	}
	return out
}
