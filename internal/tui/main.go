package tui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/lipgloss"
)

const (
	colourRed = "#ff0000"
	colourGreen = "#00ff00"
)

type Tui struct {
	std io.Writer
	err io.Writer

	errStyle lipgloss.Style
	successStyle lipgloss.Style
	infoStyle lipgloss.Style
}

func (t *Tui) Error(msg... string) {
	fmt.Fprintln(t.err, t.errStyle.Render(msg...))
}

func (t *Tui) Success(msg... string) {
	fmt.Fprintln(t.std, t.successStyle.Render(msg...))
}

func (t *Tui) Info(msg... string) {
	fmt.Fprintln(t.std, t.infoStyle.Render(msg...))
}

func NewTui(std io.Writer, err io.Writer) *Tui {
	return &Tui{
		std: std,
		err: err,
		errStyle: lipgloss.NewStyle().Foreground(lipgloss.Color(colourRed)),
		infoStyle: lipgloss.NewStyle(),
		successStyle: lipgloss.NewStyle().Foreground(lipgloss.Color(colourGreen)),
	}
}