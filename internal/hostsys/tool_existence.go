package hostsys

import (
	"fmt"
	"os/exec"

	"github.com/panoptescloud/orca/internal/common"
)

type CheckMethod string

const (
	CheckMethodPathLookup CheckMethod = "path_lookup"
	CheckMethodSuccessfulExitCode CheckMethod = "successful_exit_code"
)

type Tool struct {
	Cmd string
	Method CheckMethod
	SubCmds []Tool
}

func RequiredTools() []Tool {
	return []Tool{
		{
			Cmd: "docker",
			Method: CheckMethodPathLookup,
			SubCmds: []Tool{
				{
					Cmd: "compose",
					Method: CheckMethodSuccessfulExitCode,
				},
			},
		},
		{
			Cmd: "git",
			Method: CheckMethodPathLookup,
		},
	}
}

type EnsureToolsExistDTO struct {
	Tools []Tool
}

func buildFullCommand(tools... Tool) []string {
	parts := make([]string, len(tools))

	for i, t := range tools {
		parts[i] = t.Cmd
	}

	return parts
}

func ensureToolExists(tool Tool, parents []Tool) error {

	switch tool.Method {
	case CheckMethodPathLookup:
		if _, err := exec.LookPath(tool.Cmd); err != nil {
			return common.ErrToolsNotFoundOnSystem{
				Errs: []common.ErrToolNotFoundOnSystem{
					{
						Tool: tool.Cmd,
					},
				},
			}
		}
	case CheckMethodSuccessfulExitCode:
		cmdParts := buildFullCommand(append(parents, tool)...)
		cmd := exec.Command(cmdParts[0], cmdParts[1:]...)

		if err := cmd.Start(); err != nil {
			return err
		}

		if err := cmd.Wait(); err != nil {
			if exiterr, ok := err.(*exec.ExitError); ok {
				if exiterr.ExitCode() != 0 {
					return common.ErrToolsNotFoundOnSystem{
						Errs: []common.ErrToolNotFoundOnSystem{
							{
								Tool: tool.Cmd,
							},
						},
					}
				}
			}

			return err
		}
	default:
		return common.ErrInvalidStrategy{
			Msg: fmt.Sprintf("'%s' is not a valid CheckMethod for cmd '%s'", tool.Method, tool.Cmd),
		}
	}

	if len(tool.SubCmds) == 0 {
		return nil
	}

	parents = append(parents, tool)

	subErrs := []common.ErrToolNotFoundOnSystem{}

	for _, sub := range tool.SubCmds {
		if err := ensureToolExists(sub, parents); err != nil {
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

func EnsureToolsExist(dto EnsureToolsExistDTO) error {
	errs := []common.ErrToolNotFoundOnSystem{}

	for _, t := range dto.Tools {
		err := ensureToolExists(t, []Tool{})

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

