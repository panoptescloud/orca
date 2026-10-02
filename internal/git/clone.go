package git

import (
	"github.com/adamkirk/orca/internal/common"
	"github.com/adamkirk/orca/internal/hostsys"
)

func (g *Git) Clone(repoURL string, target string) error {
	if repoURL == "" || target == "" {
		return common.ErrInvalidInput{
			To:  "git.clone",
			Msg: "'repoURL' and 'target' cannot be empty",
		}
	}
	return g.exec.Exec(
		"git",
		[]string{
			"clone",
			repoURL,
			target,
		},
		hostsys.WithHostIO(),
	)
}
