package git

import (
	"github.com/panoptescloud/orca/internal/common"
)

type CheckoutDTO struct {
	Name string
}

func (g *Git) performCheckout(branch string) error {
	_, _, err := g.exec.Exec("git", []string{
		"checkout",
		branch,
	})

	return err
}

func (g *Git) Checkout(dto CheckoutDTO) error {
	if err := g.mustBeInAGitRepository(); err != nil {
		return err
	}

	branches, err := g.searchBranches(SearchBranchesDTO{
		Search: dto.Name,
	})

	if err != nil {
		return err
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
