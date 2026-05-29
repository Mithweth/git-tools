package main

import (
	"context"
	"fmt"
	"github.com/Mithweth/git-tools/internal/domain"
	"github.com/Mithweth/git-tools/internal/git"
	"github.com/Mithweth/git-tools/internal/providers"
	"os"
)

func CreatePR() (string, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var err error
	baseBranch := os.Getenv("BASE_BRANCH")
	if baseBranch == "" {
		baseBranch, err = git.GetDefaultBranch()
		if err != nil {
			return "", err
		}
	}
	headBranch := os.Getenv("BRANCH_NAME")
	if headBranch == "" {
		headBranch, err = git.GetCurrentBranch()
		if err != nil {
			return "", err
		}
	}
	message := os.Getenv("COMMIT_MESSAGE")
	if message == "" {
		message, err = git.GetLastCommitMessage()
		if err != nil {
			return "", err
		}
	}
	provider, owner, repoName, err := git.GetRepositoryName()
	if err != nil {
		return "", err
	}
	var token string
	switch provider {
	case domain.ProviderGitHub:
		token = os.Getenv("GITHUB_TOKEN")
	case domain.ProviderGitLab:
		token = os.Getenv("GITLAB_TOKEN")
	}
	if token == "" {
		return "", fmt.Errorf("token is not set")
	}
	repo, err := providers.NewRepository(ctx, provider, token, owner, repoName)
	if err != nil {
		return "", err
	}

	url, err := repo.CreatePullRequest(ctx, message, headBranch, baseBranch)
	if err != nil {
		return "", err
	}
	return url, nil
}

func main() {
	url, err := CreatePR()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(url)
	os.Exit(0)
}
