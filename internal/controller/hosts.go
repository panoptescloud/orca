package controller

type HostsDTO struct {
	Workspace string
}

func (c *Controller) Hosts(dto HostsDTO) error {
	ctx, err := c.contextResolver.Resolve(dto.Workspace, "")

	if err != nil {
		return err
	}

	ws, err := c.workspaceRepo.Load(ctx.Workspace.Name)

	if err != nil {
		return err
	}

	if err := c.etcHostsManager.SyncForWorkspace(ws); err != nil {
		return c.tui.RecordIfError("Failed to save /etc/hosts file!", err)
	}

	c.tui.Success("Hosts are now added to the /etc/hosts file!")

	return nil
}
