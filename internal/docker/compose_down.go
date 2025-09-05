package docker

import (
	"fmt"

	"github.com/panoptescloud/orca/internal/common"
	"github.com/panoptescloud/orca/internal/hostsys"
)

func buildComposeDownCommand(ws *common.Workspace, p *common.Project) []string {
	return []string{
		"docker",
		"compose",
		"-f",
		p.Config.ComposeFiles.Primary,
		"-p",
		fmt.Sprintf("orca-%s-%s", ws.Name, p.Name),
		"down",
	}
}

func (c *Compose) Down(ws *common.Workspace, p *common.Project) error {
	c.tui.Info(fmt.Sprintf("%s:%s[%s] stopping...", p.Name, ws.Name, p.ProjectDir))
	c.tui.NewLine()
	cmd := buildComposeDownCommand(ws, p)

	err := c.cli.Exec(cmd[0], cmd[1:], hostsys.WithHostIO(), hostsys.ChdirOpt(p.ProjectDir))

	if err != nil {
		return c.tui.RecordIfError(fmt.Sprintf("%s:%s[%s] failed!", p.Name, ws.Name, p.ProjectDir), err)
	}

	c.tui.Success(fmt.Sprintf("%s:%s[%s] stopped!", p.Name, ws.Name, p.ProjectDir))
	c.tui.NewLine()

	return nil
}
