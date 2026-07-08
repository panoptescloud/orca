package controller

import (
	"github.com/panoptescloud/orca/internal/common"
	"github.com/panoptescloud/orca/pkg/dag"
)

type UpDTO struct {
	Workspace string
	Project   string
}

func determineStartupOrder(ctx common.ExecutionContext) ([]string, error) {
	if ctx.Project != nil {
		return []string{
			ctx.Project.Name,
		}, nil
	}

	g, err := dag.NewGraph(ctx.Workspace.Projects)

	if err != nil {
		return nil, err
	}

	return g.TopologicalKeysFromRoots()
}

func (c *Controller) startServices(ctx common.ExecutionContext) error {
	ordered, err := determineStartupOrder(ctx)

	if err != nil {
		return err
	}

	for _, p := range ordered {
		// This shouldn't ever really return an error at this point, if it does
		// something went wrong while building the runtime context
		projectConfig, err := ctx.Workspace.GetProject(p)

		if err != nil {
			return err
		}

		if err := c.compose.Up(ctx.Workspace, projectConfig); err != nil {
			return err
		}
	}

	return nil
}

func (c *Controller) Up(dto UpDTO) error {
	ctx, err := c.contextResolver.Resolve(dto.Workspace, dto.Project)

	if err != nil {
		return err
	}

	return c.startServices(ctx)
}
