package commands

import "fmt"

type ErrorInvalidGitTreeOutput struct{}

func (e ErrorInvalidGitTreeOutput) Error() string { return "invalid git tree output format" }

type ErrorInvalidGitLogOutput struct{}

func (e ErrorInvalidGitLogOutput) Error() string { return "invalid git log output format" }

type ErrorInvalidGitBlameOutput struct {
	Info string
}

func (e ErrorInvalidGitBlameOutput) Error() string {
	return fmt.Sprintf("invalid git blame output format (%s)", e.Info)
}

type ErrorCommandExecution struct {
	C string
	E error
}

func (e ErrorCommandExecution) Error() string {
	return fmt.Sprintf("command %q failed (%v)", e.C, e.E)
}
