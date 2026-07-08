package controller

type ProvisionDTO struct {
	Workspace string
	Project   string
}

func (c *Controller) Provision(dto ProvisionDTO) error {
	ctx, err := c.contextResolver.Resolve(dto.Workspace, dto.Project)

	if err != nil {
		return err
	}

	return c.provisionerRunner.RunAll(ctx.Workspace, ctx.Project)
}
