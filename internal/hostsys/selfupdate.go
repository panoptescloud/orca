package hostsys

import (
	"bytes"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/minio/selfupdate"
	"github.com/panoptescloud/orca/internal/common"
)

type VersioningStrategy string

const VersioningStrategyLatest VersioningStrategy = "latest"
const VersioningStrategySpecific VersioningStrategy = "specific"

func (self *HostSystem) replaceSelf(version string) error {
	os, err := self.getOrcaOSForDownload()

	if err != nil {
		return err
	}

	arch, err := self.getOrcaArchForDownload()

	if err != nil {
		return err
	}

	url := fmt.Sprintf(
		"https://github.com/%s/%s/releases/download/v%s/orca_%s_%s.tar.gz",
		orcaRepoOrg,
		orcaRepoName,
		version,
		os,
		arch,
	)

	contents, err := self.getExecutableFromArchiveUrl(url, "orca")

	if err != nil {
		return err
	}

	reader := bytes.NewReader(contents)

	return selfupdate.Apply(reader, selfupdate.Options{})
}

func (s *HostSystem) ensureSelfVersionExists(v *semver.Version) error {
	exists, err := s.github.VersionExists(orcaRepoOrg, orcaRepoName, fmt.Sprintf("v%s", v.String()))

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

func (s *HostSystem) getVersionForSelf(dto UpdateSelfDTO) (*semver.Version, error) {
	if dto.Strategy == VersioningStrategySpecific {
		v, err := semver.NewVersion(dto.SpecifiedVersion)

		if err != nil {
			return nil, err
		}

		if err := s.ensureSelfVersionExists(v); err != nil {
			return nil, err
		}

		return v, nil
	}

	return s.github.LatestVersion(orcaRepoOrg, orcaRepoName)
}

type UpdateSelfDTO struct {
	Strategy         VersioningStrategy
	SpecifiedVersion string
}

func (self *HostSystem) UpdateSelf(dto UpdateSelfDTO) error {
	v, err := self.getVersionForSelf(dto)

	if err != nil {
		return self.tui.RecordIfError("Failed to find version!", err)
	}

	self.tui.Info(fmt.Sprintf("Switching to version: %s", v.String()))

	err = self.replaceSelf(v.String())

	if err != nil {
		return self.tui.RecordIfError(
			"Failed to apply update!",
			err,
		)
	}

	self.tui.Success("Successfully switched!")
	return nil
}
