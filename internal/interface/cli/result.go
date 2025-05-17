package cli

import "github.com/yusadeol/go-commit/internal/Domain/vo"

type Result struct {
	ExitCode ExitCode
	Message  *vo.ColoredText
}

type ExitCode int

const (
	ExitCodeSuccess           ExitCode = 0
	ExitCodeError             ExitCode = 1
	ExitCodeInvalidUsage      ExitCode = 2
	ExitCodeCommandNotFound   ExitCode = 127
	ExitCodePermissionDenied  ExitCode = 126
	ExitCodeInterruptedByUser ExitCode = 130
	ExitCodeOutOfMemory       ExitCode = 137
)

func NewResult() *Result {
	return &Result{ExitCode: ExitCodeSuccess}
}
