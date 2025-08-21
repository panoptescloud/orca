package updater

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"runtime"

	"github.com/Masterminds/semver/v3"
	"github.com/minio/selfupdate"
	"github.com/panoptescloud/orca/internal/common"
)

const githubOwner = "panoptescloud"
const githubRepo = "orca"

type VersioningStrategy string

const VersioningStrategyLatest VersioningStrategy = "latest"
const VersioningStrategySpecific VersioningStrategy = "specific"

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
	github githubClient
	tui    tui
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

func (s *SelfUpdater) getOSForGithubReleaseURL() (string, error) {
	switch os := runtime.GOOS; os {
	case "darwin":
		return "Darwin", nil
	case "linux":
		return "Linux", nil
	case "windows":
		return "Windows", nil
	default:
		return "", common.ErrUnsupportedOS{
			OS: os,
		}
	}
}

func (s *SelfUpdater) getArchForGithubReleaseURL() (string, error) {
	switch arch := runtime.GOARCH; arch {
	case "amd64":
		return "X86_64", nil
	default:
		return "", common.ErrUnsupportedArchitecture{
			Arch: arch,
		}
	}
}

func (s *SelfUpdater) extractExecutableFromArchive(executable []byte) (io.Reader, error) {
	buf := bytes.NewReader(executable)
	gzr, err := gzip.NewReader(buf)
	if err != nil {
		return nil, err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()

		switch {

		// if no more files are found return
		case err == io.EOF:
			return nil, common.ErrInvalidArchive{}

		// return any other error
		case err != nil:
			return nil, err

		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		if header.Name == "orca" {
			return tr, nil
		}
	}
}

func (s *SelfUpdater) getExecutable(v *semver.Version) (io.Reader, error) {
	os, err := s.getOSForGithubReleaseURL()

	if err != nil {
		return nil, err
	}

	arch, err := s.getArchForGithubReleaseURL()

	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf(
		"https://github.com/%s/%s/releases/download/v%s/orca_%s_%s.tar.gz",
		githubOwner,
		githubRepo,
		v.String(),
		os,
		arch,
	)

	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	executable, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	if err := resp.Body.Close(); err != nil {
		return nil, err
	}

	binary, err := s.extractExecutableFromArchive(executable)

	if err != nil {
		return nil, err
	}

	return binary, nil
}

type ApplyDTO struct {
	Strategy         VersioningStrategy
	SpecifiedVersion string
}

func (s *SelfUpdater) Apply(dto ApplyDTO) error {
	v, err := s.getVersion(dto)

	if err != nil {
		return s.tui.RecordIfError("Failed to find version!", err)
	}

	s.tui.Info(fmt.Sprintf("Switching to version: %s", v.String()))

	binary, err := s.getExecutable(v)

	if err != nil {
		return s.tui.RecordIfError("Failed to get executable contents!", err)
	}

	err = selfupdate.Apply(binary, selfupdate.Options{})
	if err != nil {
		s.tui.RecordIfError("Failed to apply update!", err)
	}

	s.tui.Success(fmt.Sprintf("Successfully switched to version: %s", v.String()))

	return nil
}

func NewSelfUpdater(github githubClient, tui tui) *SelfUpdater {
	return &SelfUpdater{
		github: github,
		tui:    tui,
	}
}
