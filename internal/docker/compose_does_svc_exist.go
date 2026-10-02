package docker

import (
	"github.com/panoptescloud/orca/internal/common"
)

// TODO: guard against nil arguments
func (c *Compose) DoesSvcExist(ws *common.Workspace, p *common.Project, service string) (bool, error) {
	spec, err := c.LoadSpec(ws, p)

	if err != nil {
		c.tui.RecordIfError("Failed to load docker compose spec for project", err)
		return false, err
	}

	_, ok := spec.Services[service]

	return ok, nil
}
