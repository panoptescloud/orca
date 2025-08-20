package git

import "github.com/panoptescloud/orca/internal/hostsys"

type PushDTO struct {
	Force bool
}

func (g *Git) Push(dto PushDTO) error {
	if err := g.mustBeInAGitRepository(); err != nil {
		return err
	}

	branch, err := g.GetCurrentBranch()

	if err != nil {
		return g.tui.RecordIfError("Failed to determine current branch, and no branch was supplied!", err)
	}

	args := []string{
		"push",
	}

	if dto.Force {
		args = append(args, "-f")
	}

	args = append(args, "origin", branch)

	return g.exec.Exec("git", args, hostsys.WithHostIO())
}
