package controller

type LogsDTO struct {
	Workspace string
	Project   string
	Service   string
}

func (c *Controller) Logs(dto LogsDTO) error {
	ctx, err := c.resolveContext(dto.Workspace, dto.Project)

	if err != nil {
		return err
	}

	return c.compose.Logs(ctx.Workspace, ctx.Project, dto.Service)
}
