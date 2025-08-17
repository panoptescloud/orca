package hostsys

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/panoptescloud/orca/internal/common"
)

type Tool struct {
	Cmd string
	SubCmds []Tool
}

func defaultRequirements() []Tool {
	return []Tool{
		{
			Cmd: "docker",
			SubCmds: []Tool{
				{
					Cmd: "compose",
				},
			},
		},
		{
			Cmd: "git",
		},
	}
}

func buildFullCommand(tools... Tool) []string {
	parts := make([]string, len(tools))

	for i, t := range tools {
		parts[i] = t.Cmd
	}

	return parts
}

func (hs *HostSystem) checkByPathLookup(tool Tool) error {
	if _, err := exec.LookPath(tool.Cmd); err != nil {
		hs.tui.Error(fmt.Sprintf("'%s' was not found in the PATH", tool.Cmd))

		return common.ErrToolsNotFoundOnSystem{
			Errs: []common.ErrToolNotFoundOnSystem{
				{
					Tool: tool.Cmd,
				},
			},
		}
	}

	hs.tui.Success(fmt.Sprintf("'%s' was found in the PATH!", tool.Cmd))

	return nil
}

func (hs *HostSystem) checkByExitCode(tool Tool, parents []Tool) error {
	cmdParts := buildFullCommand(append(parents, tool)...)
	cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
	fullCmdString := strings.Join(cmdParts, " ")

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if exiterr.ExitCode() != 0 {
				hs.tui.Error(fmt.Sprintf("'%s' returned a non-zero exit code!", fullCmdString))

				return common.ErrToolsNotFoundOnSystem{
					Errs: []common.ErrToolNotFoundOnSystem{
						{
							Tool: tool.Cmd,
						},
					},
				}
			}
		}

		hs.tui.Error(fmt.Sprintf("Something wen wrong while checking '%s', this is most likely a bug! (whoopsie)", fullCmdString))
		return err
	}

	hs.tui.Success(fmt.Sprintf("'%s' could be executed!", fullCmdString))

	return nil
}

func (hs *HostSystem) ensureToolExists(tool Tool, parents []Tool) error {
	if len(parents) == 0 {
		if err := hs.checkByPathLookup(tool); err != nil {
			return err
		}
	} else {
		if err := hs.checkByExitCode(tool, parents); err != nil {
			return err
		}
	}

	if len(tool.SubCmds) == 0 {
		return nil
	}

	parents = append(parents, tool)

	subErrs := []common.ErrToolNotFoundOnSystem{}

	for _, sub := range tool.SubCmds {
		if err := hs.ensureToolExists(sub, parents); err != nil {
			if notFoundErr, ok := err.(common.ErrToolsNotFoundOnSystem); ok {
				subErrs = append(subErrs, notFoundErr.Errs...)
			}

			return err
		}
	}

	if len(subErrs) > 0 {
		return common.ErrToolsNotFoundOnSystem{
			Errs: subErrs,
		}
	}

	return nil
}

func (hs *HostSystem) VerifySetup() error {
	errs := []common.ErrToolNotFoundOnSystem{}

	for _, t := range defaultRequirements() {
		err := hs.ensureToolExists(t, []Tool{})

		if err != nil {
			if toolNotFoundErr, ok := err.(common.ErrToolsNotFoundOnSystem); ok {
				errs = append(errs, toolNotFoundErr.Errs...)
				continue
			}

			return err
		}
	}

	if len(errs) > 0 {
		return common.ErrToolsNotFoundOnSystem{
			Errs: errs,
		}
	}

	return nil
}

