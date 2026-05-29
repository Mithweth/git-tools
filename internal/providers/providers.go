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

func New(ctx context.Context, r *domain.Repository, token string) (Repository, error) {
	switch r.Provider {
	case domain.ProviderGitHub:
		return github.Repository(ctx, token, r.Owner, r.Name), nil
	case domain.ProviderGitLab:
		return gitlab.Repository(ctx, token, r.Owner, r.Name), nil
	default:
		return nil, fmt.Errorf("unsupported provider: %s", r.Provider)
	}
}
