package git

import (
	"fmt"

	"github.com/panoptescloud/orca/internal/common"
	"github.com/panoptescloud/orca/internal/hostsys"
)

type UndoLastXCommitsDTO struct {
	// The amount of commits to remoe
	Amount           int
	SkipConfirmation bool
}

func (g *Git) UndoLastXCommits(dto UndoLastXCommitsDTO) error {
	if err := g.mustBeInAGitRepository(); err != nil {
		return err
	}

	if !dto.SkipConfirmation {
		result, err := g.tui.PresentChoices([]string{
			ConfirmNo,
			ConfirmYes,
		}, "You might lose some precious changes, are you sure?")

		if err != nil {
			if _, ok := err.(common.ErrUserAbortedExecution); ok {
				return err
			}

			return g.tui.RecordIfError("Something went wrong, this is most likely a bug!", err)
		}

		if result == ConfirmNo {
			g.tui.Info("Aborting")
			return nil
		}
	}

	return g.exec.Exec("git", []string{
		"reset",
		"--hard",
		fmt.Sprintf("HEAD~%d", dto.Amount),
	}, hostsys.WithHostIO())
}
