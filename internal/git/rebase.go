package git

import (
	"fmt"

	"github.com/adamkirk/orca/internal/hostsys"
)

type RebaseInteractivelyDTO struct {
	// The amount of commits to include in the rebase
	Amount int
	// The directory to run the command in, defaults to the current directory
	Dir string
}

func (g *Git) RebaseInteractively(dto RebaseInteractivelyDTO) error {
	if err := g.mustBeInAGitRepository(dto.Dir); err != nil {
		return err
	}

	return g.exec.Exec("git", []string{
		"rebase",
		"-i",
		fmt.Sprintf("HEAD~%d", dto.Amount),
	}, withDir(dto.Dir, hostsys.WithHostIO())...)
}
