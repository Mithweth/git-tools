package git

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"strings"
)

func GetDefaultBranch() (string, error) {
	repo, err := git.PlainOpen(".")
	if err != nil {
		return "", err
	}
	ref, err := repo.Reference(plumbing.ReferenceName("refs/remotes/origin/HEAD"), false)
	if err != nil {
		if strings.Contains(err.Error(), "reference not found") {
			return "", fmt.Errorf("%w: please try to run: git remote set-head origin --auto", err)
		}
		return "", err
	}
	branch := ref.Target().Short()

	return strings.TrimPrefix(branch, "origin/"), nil
}

func GetCurrentBranch() (string, error) {
	repo, err := git.PlainOpen(".")
	if err != nil {
		return "", err
	}

	head, err := repo.Head()
	if err != nil {
		return "", err
	}

	if !head.Name().IsBranch() {
		return "", fmt.Errorf("detached HEAD at %s", head.Hash())
	}

	branch := head.Name().Short()

	return strings.TrimPrefix(branch, "origin/"), nil
}

func GetRepositoryName() (string, string, error) {
	repo, err := git.PlainOpen(".")
	if err != nil {
		return "", "", err
	}

	remote, err := repo.Remote("origin")
	if err != nil {
		return "", "", err
	}

	url := remote.Config().URLs[0]
	url = strings.TrimSuffix(url, ".git")

	if strings.HasPrefix(url, "git@github.com:") {
		url = strings.TrimPrefix(url, "git@github.com:")
	} else if strings.HasPrefix(url, "https://github.com/") {
		url = strings.TrimPrefix(url, "https://github.com/")
	} else {
		return "", "", fmt.Errorf("unsupported github url: %s", url)
	}

	parts := strings.Split(url, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid github url: %s", url)
	}

	return parts[0], parts[1], nil
}

func GetLastCommitMessage() (string, error) {
	repo, err := git.PlainOpen(".")
	if err != nil {
		return "", err
	}

	head, err := repo.Head()
	if err != nil {
		return "", err
	}

	commit, err := repo.CommitObject(head.Hash())
	if err != nil {
		return "", err
	}

	return strings.SplitN(commit.Message, "\n", 2)[0], nil
}

func CommitsSinceOriginHead(repo *git.Repository) ([]plumbing.Hash, error) {
	repo, err := git.PlainOpen(".")
	if err != nil {
		return nil, err
	}
	ref, err := repo.Reference(plumbing.ReferenceName("refs/remotes/origin/HEAD"), true)
	if err != nil {
		return nil, err
	}

	baseHash := ref.Hash()

	head, err := repo.Head()
	if err != nil {
		return nil, err
	}

	commit, err := repo.CommitObject(head.Hash())
	if err != nil {
		return nil, err
	}

	var commits []plumbing.Hash

	for {
		if commit.Hash == baseHash {
			break
		}

		commits = append(commits, commit.Hash)

		if commit.NumParents() == 0 {
			return nil, fmt.Errorf("base commit %s not found in current branch history", baseHash)
		}

		commit, err = commit.Parent(0)
		if err != nil {
			return nil, err
		}
	}

	return commits, nil
}
