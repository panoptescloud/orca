package controller

import (
	"fmt"
)

type HostsDTO struct {
	Workspace string
}

func (c *Controller) Hosts(dto HostsDTO) error {
	ws, err := c.workspaceRepo.Load(dto.Workspace)

	if err != nil {
		return err
	}

	for _, h := range ws.GetUniqueHosts() {
		c.tui.Info(fmt.Sprintf("127.0.0.1    %s", h))
	}

	return nil
}
