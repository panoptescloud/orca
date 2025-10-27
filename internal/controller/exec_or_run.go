package controller

type ExecDTO struct {
	Workspace string
	Project   string
	Service   string
	Args      []string
}

func (c *Controller) ExecOrRun(dto ExecDTO) error {
	ctx, err := c.resolveContext(dto.Workspace, dto.Project)

	if err != nil {
		return err
	}

	isRunning, err := c.compose.IsSvcRunning(ctx.Workspace, ctx.Project, dto.Service)

	if err != nil {
		return err
	}

	if !isRunning {
		return c.compose.Run(ctx.Workspace, ctx.Project, dto.Service, dto.Args)
	}

	return c.compose.Exec(ctx.Workspace, ctx.Project, dto.Service, dto.Args)
}
