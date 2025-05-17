package vo

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
