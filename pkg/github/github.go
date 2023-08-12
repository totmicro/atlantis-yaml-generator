package github

import (
	"context"
	"strconv"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type GithubRequest struct {
	AuthToken         string
	Owner             string
	Repo              string
	PullRequestNumber string
}

func newGitHubClient(authToken string) *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: authToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

func GetChangedFiles(gh GithubRequest) ([]string, error) {
	var changedFiles []string
	prNum, err := strconv.Atoi(gh.PullRequestNumber)
	if err != nil {
		return nil, err
	}
	client := newGitHubClient(gh.AuthToken)
	files, _, err := client.PullRequests.ListFiles(context.Background(), gh.Owner, gh.Repo, prNum, nil)

	if err != nil {
		return nil, err
	}
	for _, file := range files {
		changedFiles = append(changedFiles, *file.Filename)
	}
	return changedFiles, err
}
