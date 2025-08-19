package common

import (
	"fmt"
	"strings"
)

type ErrToolsNotFoundOnSystem struct {
	Errs []ErrToolNotFoundOnSystem
}

func (err ErrToolsNotFoundOnSystem) Error() string {
	tools := make([]string, len(err.Errs))

	for i, e := range err.Errs {
		tools[i] = e.Tool
	}

	return fmt.Sprintf("tools not available on PATH: %s", strings.Join(tools, ", "))
}

type ErrToolNotFoundOnSystem struct {
	Tool string
}

func (err ErrToolNotFoundOnSystem) Error() string {
	return fmt.Sprintf("tool not available on PATH: %s", err.Tool)
}

type ErrInvalidStrategy struct {
	Msg string
}

func (err ErrInvalidStrategy) Error() string {
	return fmt.Sprintf("invalid strategy: %s", err.Msg)
}

type ErrInvalidExecutionContext struct {
	Msg string
}

func (err ErrInvalidExecutionContext) Error() string {
	return fmt.Sprintf("invalid execution context: %s", err.Msg)
}

type ErrCommandExecutionFailed struct {
	Msg string
}

func (err ErrCommandExecutionFailed) Error() string {
	return fmt.Sprintf("command execution failed: %s", err.Msg)
}

type ErrUserAbortedExecution struct{}

func (err ErrUserAbortedExecution) Error() string {
	return "user aborted execution"
}

type ErrCouldNotDetermineCurrentBranch struct{}

func (err ErrCouldNotDetermineCurrentBranch) Error() string {
	return "could not determine current branch"
}

type ErrNoBranchesFound struct{}

func (err ErrNoBranchesFound) Error() string {
	return "no branches found"
}
