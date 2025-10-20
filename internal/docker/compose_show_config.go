package docker

import (
	"errors"

	"github.com/panoptescloud/orca/internal/common"
	"github.com/panoptescloud/orca/internal/hostsys"
)

// TODO: guard against nil inputs
func (c *Compose) ShowConfig(ws *common.Workspace, p *common.Project) error {
	if err := c.goToProject(p); err != nil {
		return err
	}

	overlay, err := c.getOverlay(ws, p)
	if err != nil {
		return c.tui.RecordIfError("Failed to generate overlays!", err)
	}

	cmd := buildBaseComposeCommand(ws, p, overlay)
	cmd = append(cmd, "config")

	stdErrOpt, stderr := hostsys.WithStderr()
	stdOutOpt, stdout := hostsys.WithStdout()

	err = c.cli.Exec(cmd[0], cmd[1:], stdErrOpt, stdOutOpt, hostsys.ChdirOpt(p.ProjectDir))

	if err != nil {
		return c.tui.RecordIfError("Failed to display config!", errors.New(stderr.String()))
	}

	c.tui.Info(stdout.String())

	return nil
}
