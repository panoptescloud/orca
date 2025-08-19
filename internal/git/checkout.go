package git

import (
	"github.com/panoptescloud/orca/internal/common"
	"github.com/panoptescloud/orca/internal/hostsys"
)

type CheckoutDTO struct {
	Name string
}

func (g *Git) performCheckout(branch string) error {
	err := g.exec.Exec("git", []string{
		"checkout",
		branch,
	}, hostsys.WithHostIO())

	return err
}

func (g *Git) Checkout(dto CheckoutDTO) (string, error) {
	if err := g.mustBeInAGitRepository(); err != nil {
		return "", err
	}

	// Checkout the previously checked out branch, handle it as a special case
	if dto.Name == "-" {
		err := g.performCheckout("-")

		if err != nil {
			return "", err
		}

		branch, err := g.GetCurrentBranch()

		if err != nil {
			return "", g.tui.RecordIfError("Checked out successfully, but failed to determine branch name after, this may have side-affects on further processes!", err)
		}

		return branch, nil
	}

	branches, err := g.searchBranches(SearchBranchesDTO{
		Search: dto.Name,
	})

	if err != nil {
		return "", err
	}

	if branches.GetCurrent() != nil && branches.GetCurrent().Name == dto.Name {
		g.tui.Info("Branch is already checked out!")
		return branches.GetCurrent().Name, nil
	}

	branches = branches.ExcludeCurrent()

	branchesLength := len(branches)

	if branchesLength == 0 {
		g.tui.Error("No branches were found!")
		return "", common.ErrNoBranchesFound{}
	}

	if branchesLength == 1 {
		err := g.performCheckout(branches[0].Name)

		return branches[0].Name, g.tui.RecordIfError("Failed to checkout branch", err)
	}

	chosen, err := g.tui.PresentChoices(branches.Names(), "Which branch?")

	if err != nil {
		if _, ok := err.(common.ErrUserAbortedExecution); ok {
			return "", err
		}

		return "", g.tui.RecordIfError("Something went wrong, this is most likely a bug!", err)
	}

	err = g.performCheckout(chosen)

	return chosen, g.tui.RecordIfError("Failed to checkout branch", err)
}
