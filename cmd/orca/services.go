package main

import (
	"os"

	"github.com/panoptescloud/orca/internal/hostsys"
	"github.com/panoptescloud/orca/internal/tui"
)

type services struct {
	tui *tui.Tui

	hostSystem *hostsys.HostSystem
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
