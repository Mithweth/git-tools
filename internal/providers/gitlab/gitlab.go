package gitlab

import (
	"context"
	"fmt"
	gitlab "gitlab.com/gitlab-org/api/client-go"
	"time"
)

type GitLabRepository struct {
	client     *gitlab.Client
	owner      string
	repository string
}

func Repository(ctx context.Context, token, owner, repository string) *GitLabRepository {
	client, _ := gitlab.NewClient(token)

	return &GitLabRepository{
		client:     client,
		owner:      owner,
		repository: repository,
	}
}

func (glr *GitLabRepository) CreatePullRequest(ctx context.Context, message, headBranch, baseBranch string) (string, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	projectID := glr.owner + "/" + glr.repository

	mrs, _, err := glr.client.MergeRequests.ListProjectMergeRequests(
		projectID,
		&gitlab.ListProjectMergeRequestsOptions{
			State:        gitlab.Ptr("opened"),
			SourceBranch: gitlab.Ptr(headBranch),
			TargetBranch: gitlab.Ptr(baseBranch),
		},
		gitlab.WithContext(ctxWithTimeout),
	)
	if err != nil {
		return "", err
	}

	if len(mrs) > 0 {
		return "", fmt.Errorf("merge request already exists: %s", mrs[0].WebURL)
	}

	mr, _, err := glr.client.MergeRequests.CreateMergeRequest(
		projectID,
		&gitlab.CreateMergeRequestOptions{
			Title:        gitlab.Ptr(message),
			Description:  gitlab.Ptr(message),
			SourceBranch: gitlab.Ptr(headBranch),
			TargetBranch: gitlab.Ptr(baseBranch),
		},
		gitlab.WithContext(ctxWithTimeout),
	)
	if err != nil {
		return "", err
	}

	return mr.WebURL, nil
}
