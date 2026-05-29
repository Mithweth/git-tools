package providers

import (
	"context"
	"fmt"
	"github.com/Mithweth/git-tools/internal/providers/github"
	"github.com/Mithweth/git-tools/internal/providers/gitlab"
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
