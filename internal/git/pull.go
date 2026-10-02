package git

import "github.com/adamkirk/orca/internal/common"

type PullDTO struct {
	AllProjects bool
	Workspace   string
	Project     string
}

// pullInContext pulls the current branch for the resolved project (or working
// directory), or for every project in the workspace when allProjects is set
func (g *Git) pullInContext(ctx common.ExecutionContext, allProjects bool) error {
	if !allProjects {
		if ctx.Project == nil {
			return g.PullBranch(PullBranchDTO{Dir: ctx.WorkingDirectory})
		}

		return g.PullBranch(PullBranchDTO{Dir: ctx.Project.ProjectDir})
	}

	for _, p := range ctx.Workspace.Projects {
		if err := g.PullBranch(PullBranchDTO{Dir: p.ProjectDir}); err != nil {
			return err
		}
	}

	return nil
}

func (g *Git) Pull(dto PullDTO) error {
	ctx, err := g.contextResolver.Resolve(dto.Workspace, dto.Project)

	if err != nil {
		return err
	}

	return g.pullInContext(ctx, dto.AllProjects)
}
