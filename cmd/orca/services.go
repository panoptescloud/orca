package main

import (
	"os"

	"github.com/panoptescloud/orca/internal/config"
	"github.com/panoptescloud/orca/internal/controller"
	"github.com/panoptescloud/orca/internal/docker"
	"github.com/panoptescloud/orca/internal/git"
	"github.com/panoptescloud/orca/internal/github"
	"github.com/panoptescloud/orca/internal/hostsys"
	"github.com/panoptescloud/orca/internal/repository"
	"github.com/panoptescloud/orca/internal/tls"
	"github.com/panoptescloud/orca/internal/tui"
	"github.com/panoptescloud/orca/internal/workspaces"
	"github.com/spf13/afero"
)

type services struct {
	fs afero.Fs

	config *config.Config

	tui *tui.Tui

	hostSystem *hostsys.HostSystem

	executor *hostsys.Executor

	git *git.Git

	githubClient *github.GithubClient

	workspaceManager *workspaces.Manager

	controller *controller.Controller

	workspaceRepo *repository.WorkspaceRepository

	certificateManager *tls.CertificateManager

	compose                 *docker.Compose
	composeParser           *docker.ComposeParser
	composeOverlayGenerator *docker.ComposeOverlayGenerator
}

func (s *services) GetFs() afero.Fs {
	if s.fs != nil {
		return s.fs
	}

	s.fs = afero.NewOsFs()

	return s.fs
}

func (s *services) GetConfig() *config.Config {
	if s.config != nil {
		return s.config
	}

	s.config = config.NewDefaultConfig(
		s.GetFs(),
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
		s.GetFs(),
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
		s.GetFs(),
		s.GetTui(),
		s.GetConfig(),
		s.GetGit(),
		s.GetWorkspaceRepository(),
	)

	return s.workspaceManager
}

func (s *services) GetController() *controller.Controller {
	if s.controller != nil {
		return s.controller
	}

	s.controller = controller.NewController(
		s.GetConfig(),
		s.GetWorkspaceRepository(),
		s.GetCompose(),
		s.GetTui(),
	)

	return s.controller
}

func (s *services) GetWorkspaceRepository() *repository.WorkspaceRepository {
	if s.workspaceRepo != nil {
		return s.workspaceRepo
	}

	s.workspaceRepo = repository.NewWorkspaceRepository(
		s.GetFs(),
		s.GetConfig(),
	)

	return s.workspaceRepo
}

func (s *services) GetCertificateManager() *tls.CertificateManager {
	if s.certificateManager != nil {
		return s.certificateManager
	}

	s.certificateManager = tls.NewCertificateManager(
		s.GetFs(),
		s.GetWorkspaceRepository(),
		s.GetTui(),
		getTLSDir(),
	)

	return s.certificateManager
}

func (s *services) GetCompose() *docker.Compose {
	if s.compose != nil {
		return s.compose
	}

	s.compose = docker.NewCompose(
		s.GetExecutor(),
		s.GetTui(),
		s.GetComposeOverlayGenerator(),
	)

	return s.compose
}

func (s *services) GetComposeParser() *docker.ComposeParser {
	if s.composeParser != nil {
		return s.composeParser
	}

	s.composeParser = docker.NewComposeParser()

	return s.composeParser
}

func (s *services) GetComposeOverlayGenerator() *docker.ComposeOverlayGenerator {
	if s.composeOverlayGenerator != nil {
		return s.composeOverlayGenerator
	}

	cm := s.GetCertificateManager()

	s.composeOverlayGenerator = docker.NewComposeOverlayGenerator(
		s.GetFs(),
		s.GetComposeParser(),
		getOverlayDir(),
		cm.GetCertsDir(),
	)

	return s.composeOverlayGenerator
}
