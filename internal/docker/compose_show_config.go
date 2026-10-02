package docker

import (
	"github.com/panoptescloud/orca/internal/common"
)

// TODO: guard against nil inputs
func (c *Compose) ShowConfig(ws *common.Workspace, p *common.Project) error {
	spec, err := c.LoadRawSpec(ws, p)

	if err != nil {
		return err
	}

	c.tui.Info(spec)

	return nil
}
