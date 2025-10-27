package docker

import (
	"log/slog"
	"strings"

	"github.com/panoptescloud/orca/internal/common"
	"github.com/panoptescloud/orca/internal/hostsys"
)

// TODO: guard against nil arguments
func (c *Compose) IsSvcRunning(ws *common.Workspace, p *common.Project, service string) (bool, error) {
	if err := c.goToProject(p); err != nil {
		return false, err
	}

	overlay, err := c.getOverlay(ws, p)
	if err != nil {
		return false, c.tui.RecordIfError("Failed to generate overlays!", err)
	}

	baseCmd := buildBaseComposeCommand(ws, p, overlay)
	psCmd := append(baseCmd, "ps", "-q", service)

	withStdout, outBuff := hostsys.WithStdout()
	withStderr, errBuff := hostsys.WithStderr()

	err = c.cli.Exec(psCmd[0], psCmd[1:], withStdout, withStderr, hostsys.ChdirOpt(p.ProjectDir))

	if err != nil {
		slog.Debug("stderr from docker compose ps", "stderr", errBuff.String())
		return false, c.tui.RecordIfError("Failed to check if pod is running", err)
	}

	containerID := strings.TrimSpace(outBuff.String())
	if containerID == "" {
		// Service not found or not running at all
		return false, nil
	}

	withStdout, outBuff = hostsys.WithStdout()
	withStderr, errBuff = hostsys.WithStderr()

	err = c.cli.Exec("docker", []string{"inspect", "-f", "{{.State.Running}}", containerID}, withStdout, withStderr)

	if err != nil {
		slog.Debug("stderr from docker inspect", "stderr", errBuff.String())
		return false, c.tui.RecordIfError("Failed to inspect pod", err)
	}

	state := strings.TrimSpace(outBuff.String())

	return state == "true", nil
}
