package controller

type ExecDTO struct {
	Workspace string
	Project   string
	Service   string
	Args      []string
}

func (c *Controller) Exec(dto ExecDTO) error {
	ctx, err := c.resolveContext(dto.Workspace, dto.Project)

	if err != nil {
		return err
	}

	return c.compose.Exec(ctx.Workspace, ctx.Project, dto.Service, dto.Args)
}
