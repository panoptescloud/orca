package tui

import (
	"fmt"
	"io"
	"log/slog"

	"github.com/charmbracelet/lipgloss"
)

const (
	colourRed   = "#ff0000"
	colourGreen = "#00ff00"
)

type Tui struct {
	std io.Writer
	err io.Writer

	errStyle     lipgloss.Style
	successStyle lipgloss.Style
}

func (t *Tui) NewLine() {
	fmt.Fprintln(t.std, "")
}

func (t *Tui) Error(msg ...string) {
	fmt.Fprintln(t.err, t.errStyle.Render(msg...))
}

func (t *Tui) Success(msg ...string) {
	fmt.Fprintln(t.std, t.successStyle.Render(msg...))
}

func (t *Tui) Info(msgs ...string) {
	for _, msg := range msgs {
		fmt.Fprintln(t.std, msg)
	}
}

func (t *Tui) RecordIfError(msg string, err error) error {
	if err == nil {
		return nil
	}

	t.Error(msg)

	slog.Error(err.Error())

	return err
}

func NewTui(std io.Writer, err io.Writer) *Tui {
	return &Tui{
		std:          std,
		err:          err,
		errStyle:     lipgloss.NewStyle().Foreground(lipgloss.Color(colourRed)),
		successStyle: lipgloss.NewStyle().Foreground(lipgloss.Color(colourGreen)),
	}
}
