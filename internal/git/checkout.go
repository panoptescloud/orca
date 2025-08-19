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

func (g *Git) Checkout(dto CheckoutDTO) error {
	if err := g.mustBeInAGitRepository(); err != nil {
		return err
	}

	// Checkout the previously checked out branch, handle it as a special case
	if dto.Name == "-" {
		return g.performCheckout("-")
	}

	branches, err := g.searchBranches(SearchBranchesDTO{
		Search: dto.Name,
	})

	if err != nil {
		return err
	}

	if branches.GetCurrent() != nil && branches.GetCurrent().Name == dto.Name {
		g.tui.Info("Branch is already checked out!")
		return nil
	}

	branches = branches.ExcludeCurrent()

	branchesLength := len(branches)

	if branchesLength == 0 {
		g.tui.Error("No branches were found!")
		return nil
	}

	if branchesLength == 1 {
		err := g.performCheckout(branches[0].Name)

		return g.tui.RecordIfError("Failed to checkout branch", err)
	}

	chosen, err := g.tui.PresentChoices(branches.Names(), "Which branch?")

	if err != nil {
		if _, ok := err.(common.ErrUserAbortedExecution); ok {
			return nil
		}

		return g.tui.RecordIfError("Something went wrong, this is most likely a bug!", err)
	}

	err = g.performCheckout(chosen)

	return g.tui.RecordIfError("Failed to checkout branch", err)
}
