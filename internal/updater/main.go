package updater

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/panoptescloud/orca/internal/common"
)

const githubOwner = "panoptescloud"
const githubRepo = "orca"

type VersioningStrategy string

const VersioningStrategyLatest VersioningStrategy = "latest"
const VersioningStrategySpecific VersioningStrategy = "specific"

type installer interface {
	ReplaceSelf(version string) error
}

type tui interface {
	Info(msg ...string)
	Error(msg ...string)
	Success(msg ...string)
	RecordIfError(msg string, err error) error
}

type githubClient interface {
	VersionExists(owner string, repo string, v string) (bool, error)
	LatestVersion(owner string, repo string) (*semver.Version, error)
}

type SelfUpdater struct {
	github    githubClient
	tui       tui
	installer installer
}

func (s *SelfUpdater) ensureVersionExists(v *semver.Version) error {
	exists, err := s.github.VersionExists(githubOwner, githubRepo, fmt.Sprintf("v%s", v.String()))

	if err != nil {
		return err
	}

	if !exists {
		return common.ErrVersionNotFound{
			Version: v.String(),
		}
	}

	return nil
}

func (s *SelfUpdater) getVersion(dto ApplyDTO) (*semver.Version, error) {
	if dto.Strategy == VersioningStrategySpecific {
		v, err := semver.NewVersion(dto.SpecifiedVersion)

		if err != nil {
			return nil, err
		}

		if err := s.ensureVersionExists(v); err != nil {
			return nil, err
		}

		return v, nil
	}

	return s.github.LatestVersion(githubOwner, githubRepo)
}

type ApplyDTO struct {
	Strategy         VersioningStrategy
	SpecifiedVersion string
}

func (self *SelfUpdater) Apply(dto ApplyDTO) error {
	v, err := self.getVersion(dto)

	if err != nil {
		return self.tui.RecordIfError("Failed to find version!", err)
	}

	self.tui.Info(fmt.Sprintf("Switching to version: %s", v.String()))

	err = self.installer.ReplaceSelf(v.String())

	if err != nil {
		return self.tui.RecordIfError(
			"Failed to apply update!",
			err,
		)
	}

	self.tui.Success("Successfully switched!")
	return nil
}

func NewSelfUpdater(github githubClient, tui tui, installer installer) *SelfUpdater {
	return &SelfUpdater{
		github:    github,
		tui:       tui,
		installer: installer,
	}
}
