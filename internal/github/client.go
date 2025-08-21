package github

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Masterminds/semver/v3"
	"github.com/google/go-github/v74/github"
	"github.com/panoptescloud/orca/internal/common"
)

type GithubClient struct {
	client *github.Client
}

func (self *GithubClient) LatestVersion(owner string, repo string) (*semver.Version, error) {
	rel, resp, err := self.client.Repositories.GetLatestRelease(context.TODO(), owner, repo)

	if err != nil {
		return nil, common.ErrUnexpectedApiError{
			Msg: fmt.Sprintf("error from client: %s", err.Error()),
		}
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return semver.NewVersion(*rel.Name)
	case http.StatusTooManyRequests:
		return nil, common.ErrRateLimitedByApi{
			Api: "github",
		}
	default:
		return nil, common.ErrUnexpectedApiError{
			Msg: fmt.Sprintf("got %d response when getting release", resp.StatusCode),
		}
	}
}

func (self *GithubClient) VersionExists(owner string, repo string, v string) (bool, error) {
	_, resp, err := self.client.Repositories.GetReleaseByTag(context.TODO(), owner, repo, v)

	if err != nil {
		return false, common.ErrUnexpectedApiError{
			Msg: fmt.Sprintf("error from client: %s", err.Error()),
		}
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return true, nil
	case http.StatusNotFound:
		return false, nil
	case http.StatusTooManyRequests:
		return false, common.ErrRateLimitedByApi{
			Api: "github",
		}
	default:
		return false, common.ErrUnexpectedApiError{
			Msg: fmt.Sprintf("got %d response when getting release", resp.StatusCode),
		}
	}
}

func NewGithubClient() *GithubClient {
	return &GithubClient{
		client: github.NewClient(nil),
	}
}
