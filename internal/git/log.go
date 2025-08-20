package git

import (
	"strings"

	"github.com/panoptescloud/orca/internal/hostsys"
)

type LoglDTO struct {
	// The amount of commits to show
	Amount int
}

func (g *Git) Logl(dto LoglDTO) error {
	if err := g.mustBeInAGitRepository(); err != nil {
		return err
	}

	opt, stdout := hostsys.WithStdout()
	err := g.exec.Exec("git", []string{
		"log",
		"--oneline",
	}, opt)

	if err != nil {
		return g.tui.RecordIfError("Something went wrong, this is most likely a bug!", err)
	}

	linesOutput := stdout.String()

	lines := strings.Split(linesOutput, "\n")

	lineCount := min(dto.Amount, len(lines))

	output := strings.Join(lines[0:lineCount], "\n")

	g.tui.Info(output)

	return nil
}
