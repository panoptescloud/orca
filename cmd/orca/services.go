package main

import (
	"os"

	"github.com/panoptescloud/orca/internal/git"
	"github.com/panoptescloud/orca/internal/hostsys"
	"github.com/panoptescloud/orca/internal/tui"
)

type services struct {
	tui *tui.Tui

	hostSystem *hostsys.HostSystem

	executor *hostsys.Executor

	git *git.Git
}

func (s *services) GetTui() *tui.Tui {
	if s.tui != nil {
		return s.tui
	}

	s.tui = tui.NewTui(
		os.Stdout,
		os.Stderr,
	)

	return s.tui
}

func (s *services) GetHostSystem() *hostsys.HostSystem {
	if s.hostSystem != nil {
		return s.hostSystem
	}

	s.hostSystem = hostsys.NewHostSystem(
		s.GetTui(),
	)

	return s.hostSystem
}

func (s *services) GetExecutor() *hostsys.Executor {
	if s.executor != nil {
		return s.executor
	}

	s.executor = hostsys.NewExecutor()

	return s.executor
}

func (s *services) GetGit() *git.Git {
	if s.git != nil {
		return s.git
	}

	s.git = git.NewGit(
		s.GetExecutor(),
		s.GetTui(),
	)

	return s.git
}
