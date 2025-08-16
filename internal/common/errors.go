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
