package git

import (
	"github.com/panoptescloud/orca/internal/common"
	"github.com/panoptescloud/orca/internal/hostsys"
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
