package providers

import (
	"context"
	"fmt"
	"github.com/Mithweth/git-tools/internal/domain"
	"github.com/Mithweth/git-tools/internal/providers/github"
	"github.com/Mithweth/git-tools/internal/providers/gitlab"
)

type Repository interface {
	CreatePullRequest(ctx context.Context, message, headBranch, baseBranch string) (string, error)
}

func NewRepository(ctx context.Context, provider domain.GitProvider, token, owner, repository string) (Repository, error) {
	switch provider {
	case domain.ProviderGitHub:
		return github.Repository(ctx, token, owner, repository), nil
	case domain.ProviderGitLab:
		return gitlab.Repository(ctx, token, owner, repository), nil
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}
