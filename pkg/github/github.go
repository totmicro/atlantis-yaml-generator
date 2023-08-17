package github

import (
	"context"
	"strconv"

	"github.com/google/go-github/github"
	"github.com/totmicro/atlantis-yaml-generator/pkg/config"
	"golang.org/x/oauth2"
)

type GithubRequest struct {
	AuthToken         string
	Owner             string
	Repo              string
	PullRequestNumber string
}

// NewGitHubClient creates a new GitHub client with the provided auth token.
func newGitHubClient(authToken string) *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: authToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

// runGHRequest returns a list of changed files in a pull request.
func runGHRequest(authToken, owner, repo, pullReqNum string) ([]string, error) {
	var changedFiles []string
	prNum, err := strconv.Atoi(pullReqNum)
	if err != nil {
		return nil, err
	}
	client := newGitHubClient(authToken)
	files, _, err := client.PullRequests.ListFiles(context.Background(), owner, repo, prNum, nil)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		changedFiles = append(changedFiles, *file.Filename)
	}
	return changedFiles, err
}

// GetChangedFiles gets the parameters to call a ghrequest that returns a list of changed files.
func GetChangedFiles() (ChangedFiles []string, err error) {
	prChangedFiles, err := runGHRequest(
		config.GlobalConfig.Parameters["gh-token"],
		config.GlobalConfig.Parameters["base-repo-owner"],
		config.GlobalConfig.Parameters["base-repo-name"],
		config.GlobalConfig.Parameters["pull-num"])
	if err != nil {
		return []string{}, err
	}
	return prChangedFiles, err
}
