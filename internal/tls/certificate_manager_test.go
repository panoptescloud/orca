package tls_test

import (
	"testing"

	"github.com/panoptescloud/orca/internal/common"
	"github.com/panoptescloud/orca/internal/tls"
	tls_mocks "github.com/panoptescloud/orca/tests/mocks/tls"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func Test_CertificateManager_Generate(t *testing.T) {
	tests := []struct {
		name            string
		in              tls.GenerateDTO
		expect          error
		configureTui    func(*testing.T, *tls_mocks.MockTui)
		configureWsRepo func(*testing.T, *tls_mocks.MockWorkspaceRepo)
		reRun           bool
	}{
		{
			name: "success - no certificates",
			in: tls.GenerateDTO{
				WorkspaceName: "test",
			},
			expect: nil,
			configureTui: func(tt *testing.T, tui *tls_mocks.MockTui) {
				tui.EXPECT().Info([]string{"Generating Root certificate if required..."})
				tui.EXPECT().Info([]string{"Creating key: /some/path/key.pem"})
				tui.EXPECT().Success([]string{"Key /some/path/key.pem created successfully!"})
				tui.EXPECT().Info([]string{"Creating certificate: /some/path/cert.pem"})
				tui.EXPECT().Success([]string{"Certificate /some/path/cert.pem created successfully!"})
			},
			configureWsRepo: func(tt *testing.T, repo *tls_mocks.MockWorkspaceRepo) {
				repo.EXPECT().Load("test").Return(&common.Workspace{
					Projects: []common.Project{},
				}, nil)
			},
		},
		{
			name: "success - generate 1 certificate",
			in: tls.GenerateDTO{
				WorkspaceName: "test",
			},
			expect: nil,
			configureTui: func(tt *testing.T, tui *tls_mocks.MockTui) {
				tui.EXPECT().Info([]string{"Generating Root certificate if required..."})
				tui.EXPECT().Info([]string{"Creating key: /some/path/key.pem"})
				tui.EXPECT().Success([]string{"Key /some/path/key.pem created successfully!"})
				tui.EXPECT().Info([]string{"Creating certificate: /some/path/cert.pem"})
				tui.EXPECT().Success([]string{"Certificate /some/path/cert.pem created successfully!"})

				tui.EXPECT().NewLine()
				tui.EXPECT().Info([]string{"Generating certificate for '*.example.com'..."})
				tui.EXPECT().Info([]string{"Creating key: /some/path/certs/_.example.com.key"})
				tui.EXPECT().Success([]string{"Key /some/path/certs/_.example.com.key created successfully!"})
				tui.EXPECT().Info([]string{"Creating certificate: /some/path/certs/_.example.com.cert"})
				tui.EXPECT().Success([]string{"Certificate /some/path/certs/_.example.com.cert created successfully!"})
			},
			configureWsRepo: func(tt *testing.T, repo *tls_mocks.MockWorkspaceRepo) {
				repo.EXPECT().Load("test").Return(&common.Workspace{
					Projects: []common.Project{
						{
							Config: common.ProjectConfig{
								TLSCertificates: []string{
									"*.example.com",
								},
							},
						},
					},
				}, nil)
			},
		},

		{
			name: "success - re-run",
			in: tls.GenerateDTO{
				WorkspaceName: "test",
			},
			expect: nil,
			reRun:  true,
			configureTui: func(tt *testing.T, tui *tls_mocks.MockTui) {
				tui.EXPECT().Info([]string{"Generating Root certificate if required..."})
				tui.EXPECT().Info([]string{"Creating key: /some/path/key.pem"})
				tui.EXPECT().Success([]string{"Key /some/path/key.pem created successfully!"})
				tui.EXPECT().Info([]string{"Creating certificate: /some/path/cert.pem"})
				tui.EXPECT().Success([]string{"Certificate /some/path/cert.pem created successfully!"})

				tui.EXPECT().NewLine()
				tui.EXPECT().Info([]string{"Generating certificate for '*.example.com'..."})
				tui.EXPECT().Info([]string{"Creating key: /some/path/certs/_.example.com.key"})
				tui.EXPECT().Success([]string{"Key /some/path/certs/_.example.com.key created successfully!"})
				tui.EXPECT().Info([]string{"Creating certificate: /some/path/certs/_.example.com.cert"})
				tui.EXPECT().Success([]string{"Certificate /some/path/certs/_.example.com.cert created successfully!"})

				tui.EXPECT().Info([]string{"Key /some/path/key.pem already exists, skipping creation."})
				tui.EXPECT().RecordIfError("Failed to load key /some/path/key.pem", nil).Return(nil)
				tui.EXPECT().Info([]string{"Certificate /some/path/cert.pem already exists, skipping creation."})
				tui.EXPECT().Info([]string{"Key /some/path/certs/_.example.com.key already exists, skipping creation."})
				tui.EXPECT().RecordIfError("Failed to load key /some/path/certs/_.example.com.key", nil).Return(nil)
				tui.EXPECT().Info([]string{"Certificate /some/path/certs/_.example.com.cert already exists, skipping creation."})
				// tui.EXPECT().Info([]string{"Key /some/path/key.pem already exists, skipping creation."})
			},
			configureWsRepo: func(tt *testing.T, repo *tls_mocks.MockWorkspaceRepo) {
				repo.EXPECT().Load("test").Return(&common.Workspace{
					Projects: []common.Project{
						{
							Config: common.ProjectConfig{
								TLSCertificates: []string{
									"*.example.com",
								},
							},
						},
					},
				}, nil)
			},
		},

		{
			name: "success - re-run with multiple",
			in: tls.GenerateDTO{
				WorkspaceName: "test",
			},
			expect: nil,
			reRun:  true,
			configureTui: func(tt *testing.T, tui *tls_mocks.MockTui) {
				tui.EXPECT().Info([]string{"Generating Root certificate if required..."})
				tui.EXPECT().Info([]string{"Creating key: /some/path/key.pem"})
				tui.EXPECT().Success([]string{"Key /some/path/key.pem created successfully!"})
				tui.EXPECT().Info([]string{"Creating certificate: /some/path/cert.pem"})
				tui.EXPECT().Success([]string{"Certificate /some/path/cert.pem created successfully!"})

				tui.EXPECT().NewLine()
				tui.EXPECT().Info([]string{"Generating certificate for '*.example.com'..."})
				tui.EXPECT().Info([]string{"Creating key: /some/path/certs/_.example.com.key"})
				tui.EXPECT().Success([]string{"Key /some/path/certs/_.example.com.key created successfully!"})
				tui.EXPECT().Info([]string{"Creating certificate: /some/path/certs/_.example.com.cert"})
				tui.EXPECT().Success([]string{"Certificate /some/path/certs/_.example.com.cert created successfully!"})

				tui.EXPECT().Info([]string{"Generating certificate for 'blah.test'..."})
				tui.EXPECT().Info([]string{"Creating key: /some/path/certs/blah.test.key"})
				tui.EXPECT().Success([]string{"Key /some/path/certs/blah.test.key created successfully!"})
				tui.EXPECT().Info([]string{"Creating certificate: /some/path/certs/blah.test.cert"})
				tui.EXPECT().Success([]string{"Certificate /some/path/certs/blah.test.cert created successfully!"})

				tui.EXPECT().Info([]string{"Key /some/path/key.pem already exists, skipping creation."})
				tui.EXPECT().RecordIfError("Failed to load key /some/path/key.pem", nil).Return(nil)
				tui.EXPECT().Info([]string{"Certificate /some/path/cert.pem already exists, skipping creation."})
				tui.EXPECT().Info([]string{"Key /some/path/certs/_.example.com.key already exists, skipping creation."})
				tui.EXPECT().RecordIfError("Failed to load key /some/path/certs/_.example.com.key", nil).Return(nil)
				tui.EXPECT().Info([]string{"Certificate /some/path/certs/_.example.com.cert already exists, skipping creation."})

				tui.EXPECT().Info([]string{"Key /some/path/certs/blah.test.key already exists, skipping creation."})
				tui.EXPECT().RecordIfError("Failed to load key /some/path/certs/blah.test.key", nil).Return(nil)
				tui.EXPECT().Info([]string{"Certificate /some/path/certs/blah.test.cert already exists, skipping creation."})
			},
			configureWsRepo: func(tt *testing.T, repo *tls_mocks.MockWorkspaceRepo) {
				repo.EXPECT().Load("test").Return(&common.Workspace{
					Projects: []common.Project{
						{
							Config: common.ProjectConfig{
								TLSCertificates: []string{
									"*.example.com",
									"blah.test",
								},
							},
						},
					},
				}, nil)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			tui := tls_mocks.NewMockTui(t)
			if test.configureTui != nil {
				test.configureTui(tt, tui)
			}

			wsRepo := tls_mocks.NewMockWorkspaceRepo(t)
			if test.configureWsRepo != nil {
				test.configureWsRepo(tt, wsRepo)
			}

			fs := afero.NewMemMapFs()

			cm := tls.NewCertificateManager(fs, wsRepo, tui, "/some/path")

			assert.Equal(tt, test.expect, cm.Generate(test.in))

			if !test.reRun {
				return
			}

			assert.Equal(tt, test.expect, cm.Generate(test.in))
		})
	}
}
