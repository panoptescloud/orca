package main

import (
	"os"

	"github.com/panoptescloud/orca/internal/config"
	"github.com/panoptescloud/orca/internal/git"
	"github.com/panoptescloud/orca/internal/github"
	"github.com/panoptescloud/orca/internal/hostsys"
	"github.com/panoptescloud/orca/internal/tui"
	"github.com/panoptescloud/orca/internal/updater"
	"github.com/panoptescloud/orca/internal/workspaces"
	"github.com/spf13/afero"
)

type services struct {
	configManager *config.Manager

	tui *tui.Tui

	hostSystem *hostsys.HostSystem

	executor *hostsys.Executor

	git *git.Git

	githubClient *github.GithubClient

	selfUpdater *updater.SelfUpdater

	workspaceManager *workspaces.Manager
}

func (s *services) GetConfigManager() *config.Manager {
	if s.configManager != nil {
		return s.configManager
	}

	s.configManager = config.NewManager(
		s.GetTui(),
		afero.NewOsFs(),
		getConfigFilePath(),
	)

	return s.configManager
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

func (s *services) GetSelfUpdater() *updater.SelfUpdater {
	if s.selfUpdater != nil {
		return s.selfUpdater
	}

	s.selfUpdater = updater.NewSelfUpdater(
		s.GetGithubClient(),
		s.GetTui(),
		s.GetHostSystem(),
	)

	return s.selfUpdater
}

func (s *services) GetWorkspaceManager() *workspaces.Manager {
	if s.workspaceManager != nil {
		return s.workspaceManager
	}

	s.workspaceManager = workspaces.NewManager(
		afero.NewOsFs(),
		s.GetTui(),
		s.GetConfigManager(),
		s.GetGit(),
	)

	return s.workspaceManager
}
