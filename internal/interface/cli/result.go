package cli

import "github.com/yusadeol/go-commit/internal/Domain/vo"

type Result struct {
	ExitCode vo.ExitCode
	Message  *vo.ColoredText
}

func NewResult() *Result {
	return &Result{ExitCode: vo.ExitCodeSuccess}
}
