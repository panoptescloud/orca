package docker

import (
	"context"

	composecli "github.com/compose-spec/compose-go/v2/cli"
	"github.com/compose-spec/compose-go/v2/types"
)

type ComposeParser struct {
}

func (cp *ComposeParser) Parse(paths []string, envFiles []string) (*types.Project, error) {
	opts, err := composecli.NewProjectOptions(
		paths,
		composecli.WithEnvFiles(envFiles...),
	)

	if err != nil {
		return nil, err
	}

	if err := composecli.WithDotEnv(opts); err != nil {
		return nil, err
	}

	return opts.LoadProject(context.Background())
}

func NewComposeParser() *ComposeParser {
	return &ComposeParser{}
}
