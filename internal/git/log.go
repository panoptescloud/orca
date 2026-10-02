package git

import (
	"strings"

	"github.com/adamkirk/orca/internal/hostsys"
)

type LoglDTO struct {
	// The amount of commits to show
	Amount int
	// The directory to run the command in, defaults to the current directory
	Dir string
}

func (g *Git) Logl(dto LoglDTO) error {
	if err := g.mustBeInAGitRepository(dto.Dir); err != nil {
		return err
	}

	opt, stdout := hostsys.WithStdout()
	err := g.exec.Exec("git", []string{
		"log",
		"--oneline",
	}, withDir(dto.Dir, opt)...)

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
