package git

import (
	"fmt"
	"strings"

	"github.com/panoptescloud/orca/internal/common"
	"github.com/panoptescloud/orca/internal/hostsys"
)

type Branch struct {
	Name    string
	Current bool
}
type Branches []Branch

func (branches Branches) GetCurrent() *Branch {
	for _, b := range branches {
		if b.Current {
			return &b
		}
	}

	return nil
}

func (branches Branches) Names() []string {
	names := make([]string, len(branches))
	for i, b := range branches {
		names[i] = b.Name
	}

	return names
}

func (branches Branches) ExcludeCurrent() Branches {
	new := Branches{}

	for _, b := range branches {
		if b.Current {
			continue
		}

		new = append(new, b)
	}

	return new
}

type branchFilterFunc func(Branch) bool

func (branches Branches) Filter(f branchFilterFunc) Branches {
	new := Branches{}

	for _, b := range branches {
		if f(b) {
			new = append(new, b)
		}
	}

	return new
}

func (branches Branches) StrictSearch(search string) Branches {
	return branches.Filter(func(b Branch) bool {
		return strings.Contains(b.Name, search)
	})
}

type SearchBranchesDTO struct {
	Search string
}

func parseGitBranchOutput(output string) Branches {
	lines := strings.Split(output, "\n")
	branches := Branches{}

	foundCurrent := false

	for _, line := range lines {
		if line == "" {
			continue
		}

		isCurrent := false

		if !foundCurrent {
			isCurrent = strings.HasPrefix(line, "*")

			if isCurrent {
				line = strings.TrimPrefix(line, "*")
				foundCurrent = true
			}
		}

		line = strings.TrimLeft(line, " ")

		branches = append(branches, Branch{
			Name:    line,
			Current: isCurrent,
		})
	}

	return branches
}

func (g *Git) searchBranches(dto SearchBranchesDTO) (Branches, error) {
	if err := g.mustBeInAGitRepository(); err != nil {
		return nil, err
	}

	opt, stdout := hostsys.WithStdout()
	err := g.exec.Exec("git", []string{
		"branch",
		"-l",
	}, opt)

	if err != nil {
		return nil, common.ErrInvalidExecutionContext{
			Msg: "failed to execute 'git branch'",
		}
	}

	branches := parseGitBranchOutput(stdout.String())

	if dto.Search == "" {
		return branches, nil
	}

	return branches.StrictSearch(dto.Search), nil
}

func (g *Git) ShowBranches(dto SearchBranchesDTO) error {
	branches, err := g.searchBranches(dto)
	if err != nil {
		return g.tui.RecordIfError(
			"Failed to list branches!",
			err,
		)
	}

	if len(branches) == 0 {
		g.tui.Info("No matching branches found!")
		return nil
	}

	if b := branches.GetCurrent(); b != nil {
		g.tui.Success(fmt.Sprintf("[current] %s", b.Name))
		g.tui.NewLine()
	}

	for _, b := range branches.ExcludeCurrent() {
		g.tui.Info(fmt.Sprintf("%s", b.Name))
	}

	return nil
}

func (g *Git) GetCurrentBranch() (string, error) {
	branches, err := g.searchBranches(SearchBranchesDTO{})

	if err != nil {
		return "", err
	}

	if branches.GetCurrent() == nil {
		return "", common.ErrCouldNotDetermineCurrentBranch{}
	}

	return branches.GetCurrent().Name, nil
}

type PullBranchDTO struct {
	Name string
}

func (g *Git) PullBranch(dto PullBranchDTO) error {
	if err := g.mustBeInAGitRepository(); err != nil {
		return err
	}

	branchToPull := dto.Name

	if dto.Name == "" {
		current, err := g.GetCurrentBranch()

		if err != nil {
			return err
		}

		branchToPull = current
	}

	err := g.exec.Exec("git", []string{
		"pull",
		"origin",
		branchToPull,
	}, hostsys.WithHostIO())

	return g.tui.RecordIfError("Failed to pull branch!", err)
}
