package github

import (
	"context"
	"fmt"
	"github.com/google/go-github/v74/github"
	"golang.org/x/oauth2"
	"time"
)

type GitHubRepository struct {
	client     *github.Client
	owner      string
	repository string
}

func Repository(ctx context.Context, token, owner, repository string) *GitHubRepository {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	tc := oauth2.NewClient(ctx, ts)

	return &GitHubRepository{
		client:     github.NewClient(tc),
		owner:      owner,
		repository: repository,
	}
}

func (ghc *GitHubRepository) CreatePullRequest(ctx context.Context, message, headBranch, baseBranch string) (string, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Check if the PR already exists
	prs, _, err := ghc.client.PullRequests.List(
		ctxWithTimeout,
		ghc.owner,
		ghc.repository,
		&github.PullRequestListOptions{
			State: "open",
			Head:  ghc.owner + ":" + headBranch,
		},
	)
	if err != nil {
		return "", err
	}
	if len(prs) == 1 {
		return "", fmt.Errorf("pull request already exists: %s", *prs[0].HTMLURL)
	}

	pr := &github.NewPullRequest{
		Title: github.Ptr(message),
		Head:  github.Ptr(headBranch),
		Base:  github.Ptr(baseBranch),
		Body:  github.Ptr(message),
	}

	result, _, err := ghc.client.PullRequests.Create(
		ctxWithTimeout,
		ghc.owner,
		ghc.repository,
		pr,
	)
	if err != nil {
		return "", err
	}

	return *result.HTMLURL, nil
}

func (ghc *GitHubRepository) ListPullRequests(ctx context.Context, state *string) ([]*github.PullRequest, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	prOptions := &github.PullRequestListOptions{}

	if state != nil {
		prOptions.State = *state
	}

	results, _, err := ghc.client.PullRequests.List(
		ctxWithTimeout,
		ghc.owner,
		ghc.repository,
		prOptions,
	)
	if err != nil {
		return nil, err
	}

	return results, nil
}
