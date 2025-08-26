package main

import (
	"os"

	"github.com/panoptescloud/orca/internal/config"
	"github.com/panoptescloud/orca/internal/git"
	"github.com/panoptescloud/orca/internal/github"
	"github.com/panoptescloud/orca/internal/hostsys"
	"github.com/panoptescloud/orca/internal/tui"
	"github.com/panoptescloud/orca/internal/workspaces"
	"github.com/spf13/afero"
)

type services struct {
	config *config.Config

	tui *tui.Tui

	hostSystem *hostsys.HostSystem

	executor *hostsys.Executor

	git *git.Git

	githubClient *github.GithubClient

	workspaceManager *workspaces.Manager
}

func (s *services) GetConfig() *config.Config {
	if s.config != nil {
		return s.config
	}

	s.config = config.NewDefaultConfig(
		afero.NewOsFs(),
		getConfigFilePath(),
	)

	return s.config
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
		afero.NewOsFs(),
		s.GetGithubClient(),
		getToolsDir(),
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

func (s *services) GetGithubClient() *github.GithubClient {
	if s.githubClient != nil {
		return s.githubClient
	}

	s.githubClient = github.NewGithubClient()

	return s.githubClient
}

func (s *services) GetWorkspaceManager() *workspaces.Manager {
	if s.workspaceManager != nil {
		return s.workspaceManager
	}

	s.workspaceManager = workspaces.NewManager(
		afero.NewOsFs(),
		s.GetTui(),
		s.GetConfig(),
		s.GetGit(),
	)

	return s.workspaceManager
}
