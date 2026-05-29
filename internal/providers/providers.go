package providers

import (
	"context"
	"fmt"
	"github.com/Mithweth/git-tools/internal/providers/github"
	"github.com/Mithweth/git-tools/internal/providers/gitlab"
	"strings"
)

type GitProvider string

const (
	ProviderGitHub GitProvider = "github"
	ProviderGitLab GitProvider = "gitlab"
)

type Repository interface {
	CreatePullRequest(ctx context.Context, message, headBranch, baseBranch string) (string, error)
}

func NewRepository(ctx context.Context, provider GitProvider, token, owner, repository string) (Repository, error) {
	switch provider {
	case ProviderGitHub:
		return github.Repository(ctx, token, owner, repository), nil
	case ProviderGitLab:
		return gitlab.Repository(ctx, token, owner, repository), nil
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}

func ParseRepositoryURL(url string) (GitProvider, string, string, error) {
	var provider GitProvider

	switch {
	case strings.HasPrefix(url, "git@github.com:"):
		provider = ProviderGitHub
		url = strings.TrimPrefix(url, "git@github.com:")

	case strings.HasPrefix(url, "https://github.com/"):
		provider = ProviderGitHub
		url = strings.TrimPrefix(url, "https://github.com/")

	case strings.HasPrefix(url, "git@gitlab.com:"):
		provider = ProviderGitLab
		url = strings.TrimPrefix(url, "git@gitlab.com:")

	case strings.HasPrefix(url, "https://gitlab.com/"):
		provider = ProviderGitLab
		url = strings.TrimPrefix(url, "https://gitlab.com/")

	default:
		return "", "", "", fmt.Errorf("unsupported git remote url: %s", url)
	}

	parts := strings.Split(url, "/")
	if len(parts) < 2 {
		return "", "", "", fmt.Errorf("invalid git remote url: %s", url)
	}
	owner := strings.Join(parts[:len(parts)-1], "/")
	repository := parts[len(parts)-1]

	return provider, owner, repository, nil
}
