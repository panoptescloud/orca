package docker

import (
	"os/exec"

	"github.com/panoptescloud/orca/internal/common"
	"github.com/panoptescloud/orca/internal/hostsys"
)

// TODO: guard against nil inputs
func (c *Compose) Exec(ws *common.Workspace, p *common.Project, service string, cmdArgs []string) error {
	overlay, err := c.getOverlay(ws, p)
	if err != nil {
		return c.tui.RecordIfError("Failed to generate overlays!", err)
	}

	cmd := buildBaseComposeCommand(ws, p, overlay)
	cmd = append(cmd, "exec", "-it", service)
	cmd = append(cmd, cmdArgs...)

	err = c.cli.Exec(cmd[0], cmd[1:], hostsys.WithHostIO(), hostsys.ChdirOpt(p.ProjectDir))

	if err != nil {
		// If it's exit code 130, it's because they exited an interactive container
		// typically. Not sure if this is always true though, need to see how this
		// works across more use-cases.
		if exitError, ok := err.(*exec.ExitError); ok && exitError.ProcessState.ExitCode() == 130 {
			return nil
		}

		return c.tui.RecordIfError("Command failed!", err)
	}

	return nil
}
