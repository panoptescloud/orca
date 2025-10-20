package docker

import (
	"strings"

	"github.com/panoptescloud/orca/internal/common"
)

// TODO: guard against nil inputs
func (c *Compose) ShowCommand(ws *common.Workspace, p *common.Project) error {
	if err := c.goToProject(p); err != nil {
		return err
	}

	overlay, err := c.getOverlay(ws, p)
	if err != nil {
		return c.tui.RecordIfError("Failed to generate overlays!", err)
	}

	cmd := buildBaseComposeCommand(ws, p, overlay)

	c.tui.Info(strings.Join(cmd, " "))

	return nil
}
