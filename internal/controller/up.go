package controller

import (
	"fmt"

	"github.com/panoptescloud/orca/internal/common"
	"github.com/panoptescloud/orca/pkg/dag"
)

type UpDTO struct {
	Workspace string
	Project   string
}

func determineStartupOrder(ctx runtimeContext) ([]string, error) {
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

// TODO: Move to down.go when it's added
// func determineShutdownOrder(ctx runtimeContext) ([]string, error) {
// 	if ctx.Project != nil {
// 		return []string{
// 			ctx.Project.Name,
// 		}, nil
// 	}

// 	g, err := dag.NewGraph(ctx.WorkspaceConfig.Projects)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return g.TopologicalKeysFromLeaves()
// }

func (c *Controller) startProject(cfg *common.Project) error {
	fmt.Printf("starting:\n %#v\n\n", cfg)
	return nil
}

func (c *Controller) startServices(ctx runtimeContext) error {
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

		if err := c.startProject(projectConfig); err != nil {
			return err
		}
	}

	return nil
}

func (c *Controller) Up(dto UpDTO) error {
	ctx, err := c.resolveContext(dto.Workspace, dto.Project)

	if err != nil {
		return err
	}

	return c.startServices(ctx)
}
