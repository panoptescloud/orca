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

type ErrInvalidSemVerValue struct {
	Value string
}

func (err ErrInvalidSemVerValue) Error() string {
	return fmt.Sprintf("invalid semver value given: %s", err.Value)
}

type ErrVersionNotFound struct {
	Version string
}

func (err ErrVersionNotFound) Error() string {
	return fmt.Sprintf("version did not exist: %s", err.Version)
}

type ErrUnexpectedApiError struct {
	Msg string
}

func (err ErrUnexpectedApiError) Error() string {
	return fmt.Sprintf("unexpected API error: %s", err.Msg)
}

type ErrRateLimitedByApi struct {
	Api string
}

func (err ErrRateLimitedByApi) Error() string {
	return fmt.Sprintf("rate limit exceeded for '%s'", err.Api)
}

type ErrUnsupportedOS struct {
	OS string
}

func (err ErrUnsupportedOS) Error() string {
	return fmt.Sprintf("unsupported OS: %s", err.OS)
}

type ErrUnsupportedArchitecture struct {
	Arch string
}

func (err ErrUnsupportedArchitecture) Error() string {
	return fmt.Sprintf("unsupported architecture: %s", err.Arch)
}

type ErrInvalidArchive struct{}

func (err ErrInvalidArchive) Error() string {
	return "invalid archive used"
}

type ErrUnknownWorkspace struct {
	Workspace string
}

func (err ErrUnknownWorkspace) Error() string {
	return fmt.Sprintf("unknown workspace: %s", err.Workspace)
}

type ErrUnknownTool struct {
	Tool string
}

func (err ErrUnknownTool) Error() string {
	return fmt.Sprintf("unknown tool: %s", err.Tool)
}

type ErrArgumentRequired struct {
	Name  string
	Index int
}

func (err ErrArgumentRequired) Error() string {
	return fmt.Sprintf("the '%s' argument must be supplied at position %d", err.Name, err.Index)
}

type ErrTooManyArguments struct {
	Expected int
}

func (err ErrTooManyArguments) Error() string {
	return fmt.Sprintf("too many arguments supplied, expected only %d", err.Expected)
}

type ErrNotImplemented struct {
	Msg string
}

func (err ErrNotImplemented) Error() string {
	return fmt.Sprintf("not implemented: %s", err.Msg)
}
