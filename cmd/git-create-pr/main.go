package main

import (
	"context"
	"fmt"
	"github.com/Mithweth/git-tools/internal/git"
	"github.com/Mithweth/git-tools/internal/github"
	"os"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		fmt.Fprintf(os.Stderr, "GITHUB_TOKEN environment variable is not set\n")
		os.Exit(1)
	}
	var err error
	baseBranch := os.Getenv("BASE_BRANCH")
	if baseBranch == "" {
		baseBranch, err = git.GetDefaultBranch()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	}
	headBranch := os.Getenv("BRANCH_NAME")
	if headBranch == "" {
		headBranch, err = git.GetCurrentBranch()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	}
	message := os.Getenv("COMMIT_MESSAGE")
	if message == "" {
		message, err = git.GetLastCommitMessage()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	}
	owner, repo, err := git.GetRepositoryName()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	githubRepo := github.Repository(ctx, token, owner, repo)
	url, err := githubRepo.CreatePullRequest(ctx, message, headBranch, baseBranch)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(url)
	os.Exit(0)
}
