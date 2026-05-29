package main

import (
	"context"
	"fmt"
	"github.com/Mithweth/git-tools/internal/auth"
	"github.com/Mithweth/git-tools/internal/domain"
	"github.com/Mithweth/git-tools/internal/git"
	"github.com/Mithweth/git-tools/internal/providers"
	"github.com/spf13/pflag"
	"os"
)

func CreatePullRequest(message string) (string, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var err error
	_, baseBranch, err := git.GetDefaultBranch()
	if err != nil {
		return "", err
	}
	_, headBranch, err := git.GetCurrentBranch()
	if err != nil {
		return "", err
	}
	if message == "" {
		message, err = git.GetLastCommitMessage()
		if err != nil {
			return "", err
		}
	}
	repository, err := git.GetRepository()
	if err != nil {
		return "", err
	}
	var token string
	switch repository.Provider {
	case domain.ProviderGitHub:
		token = os.Getenv("GITHUB_TOKEN")
	case domain.ProviderGitLab:
		token = os.Getenv("GITLAB_TOKEN")
	}
	if token == "" {
		return "", fmt.Errorf("token is not set")
	}
	conn, err := providers.New(ctx, repository, token)
	if err != nil {
		return "", err
	}

	url, err := conn.CreatePullRequest(ctx, message, headBranch, baseBranch)
	if err != nil {
		return "", err
	}
	return url, nil
}

func Push(force bool) error {
	repository, err := git.GetRepository()
	if err != nil {
		return err
	}
	authMethod, err := auth.GetAuth(repository)
	if err != nil {
		return err
	}
	if err := git.Push(authMethod, force); err != nil {
		return err
	}
	return nil
}

func main() {
	var (
		commitMessage string
		push          bool
		force         bool
	)
	pflag.StringVarP(&commitMessage, "message", "m", "", "pull request title")
	pflag.BoolVarP(&push, "push", "p", false, "push before creating pull request")
	pflag.BoolVarP(&force, "force", "f", false, "use force-push instead push")
	pflag.Parse()
	if push {
		if errPush := Push(force); errPush != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", errPush)
			os.Exit(1)
		}
	}
	url, err := CreatePullRequest(commitMessage)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(url)
	os.Exit(0)
}
