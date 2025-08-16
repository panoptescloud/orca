package main

import (
	"os"

	"github.com/panoptescloud/orca/internal/tui"
)

type services struct {
	tui *tui.Tui
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