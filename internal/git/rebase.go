package git

import (
	"fmt"

	"github.com/panoptescloud/orca/internal/hostsys"
)

type RebaseInteractivelyDTO struct {
	// The amount of commits to include in the rebase
	Amount int
}

func (g *Git) RebaseInteractively(dto RebaseInteractivelyDTO) error {
	if err := g.mustBeInAGitRepository(); err != nil {
		return err
	}

	return g.exec.Exec("git", []string{
		"rebase",
		"-i",
		fmt.Sprintf("HEAD~%d", dto.Amount),
	}, hostsys.WithHostIO())
}
