package controller

import (
	"github.com/panoptescloud/orca/pkg/dag"
)

type DownDTO struct {
	Workspace string
	Project   string
}

func determineShutdownOrder(ctx runtimeContext) ([]string, error) {
	if ctx.Project != nil {
		return []string{
			ctx.Project.Name,
		}, nil
	}

	g, err := dag.NewGraph(ctx.Workspace.Projects)

	if err != nil {
		return nil, err
	}

	return g.TopologicalKeysFromLeaves()
}

func (c *Controller) stopServices(ctx runtimeContext) error {
	ordered, err := determineShutdownOrder(ctx)

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

		if err := c.compose.Down(ctx.Workspace, projectConfig); err != nil {
			return err
		}
	}

	return nil
}

func (c *Controller) Down(dto DownDTO) error {
	ctx, err := c.resolveContext(dto.Workspace, dto.Project)

	if err != nil {
		return err
	}

	return c.stopServices(ctx)
}
