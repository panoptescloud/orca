package docker

import (
	"fmt"

	"github.com/panoptescloud/orca/internal/common"
	"github.com/panoptescloud/orca/internal/hostsys"
)

// TODO: guard against nil inputs
func (c *Compose) Show(ws *common.Workspace, p *common.Project) error {
	overlay, err := c.getOverlay(ws, p)
	if err != nil {
		return err
	}

	cmd := buildBaseComposeCommand(ws, p, overlay)
	cmd = append(cmd, "config")

	c.tui.NewLine()

	err = c.cli.Exec(cmd[0], cmd[1:], hostsys.WithHostIO(), hostsys.ChdirOpt(p.ProjectDir))

	if err != nil {
		return c.tui.RecordIfError(fmt.Sprintf("%s:%s[%s] failed!", p.Name, ws.Name, p.ProjectDir), err)
	}

	c.tui.NewLine()

	return nil
}
