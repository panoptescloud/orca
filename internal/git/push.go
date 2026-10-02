package git

import "github.com/adamkirk/orca/internal/hostsys"

type PushDTO struct {
	Force bool
	// The directory to run the command in, defaults to the current directory
	Dir string
}

func (g *Git) Push(dto PushDTO) error {
	if err := g.mustBeInAGitRepository(dto.Dir); err != nil {
		return err
	}

	branch, err := g.GetCurrentBranch(dto.Dir)

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

	return g.exec.Exec("git", args, withDir(dto.Dir, hostsys.WithHostIO())...)
}
