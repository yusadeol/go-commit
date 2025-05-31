package cli

import "github.com/yusadeol/go-commit/internal/domain/vo"

type Result struct {
	ExitCode vo.ExitCode
	Message  *vo.MarkupText
}

func NewResult() *Result {
	return &Result{ExitCode: vo.ExitCodeSuccess, Message: vo.NewMarkupText("")}
}
