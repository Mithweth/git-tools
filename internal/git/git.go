package git

import (
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"strings"
	"time"
)

func ResolveRevision(rev string) (string, error) {
	repo, err := git.PlainOpen(".")
	if err != nil {
		return "", err
	}

	hash, err := repo.ResolveRevision(plumbing.Revision(rev))
	if err != nil {
		return "", err
	}
	if hash == nil {
		return "", fmt.Errorf("wrong reference: %s", rev)
	}

	return hash.String(), nil
}

func GetDefaultBranch() (string, string, error) {
	repo, err := git.PlainOpen(".")
	if err != nil {
		return "", "", err
	}
	ref, err := repo.Reference(plumbing.ReferenceName("refs/remotes/origin/HEAD"), false)
	if err != nil {
		if strings.Contains(err.Error(), "reference not found") {
			return "", "", fmt.Errorf("%w: please try to run: git remote set-head origin --auto", err)
		}
		return "", "", err
	}
	def := ref.Target()
	defRef, err := repo.Reference(def, true)
	if err != nil {
		return "", "", err
	}
	return defRef.Hash().String(), strings.TrimPrefix(def.Short(), "origin/"), nil
}

func GetCurrentBranch() (string, string, error) {
	repo, err := git.PlainOpen(".")
	if err != nil {
		return "", "", err
	}

	head, err := repo.Head()
	if err != nil {
		return "", "", err
	}

	if !head.Name().IsBranch() {
		return "", "", fmt.Errorf("detached HEAD at %s", head.Hash())
	}

	return head.Hash().String(), strings.TrimPrefix(head.Name().Short(), "origin/"), nil
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

func getCommitsBetween(from, until string) ([]*object.Commit, error) {
	repo, err := git.PlainOpen(".")
	if err != nil {
		return nil, err
	}
	fromHash := plumbing.NewHash(from)
	untilHash := plumbing.NewHash(until)

	commit, err := repo.CommitObject(untilHash)
	if err != nil {
		return nil, err
	}

	var commits []*object.Commit

	for {
		//nolint:staticcheck
		if commit.Hash == fromHash {
			break
		}
		commits = append(commits, commit)
		if commit.NumParents() == 0 {
			return nil, fmt.Errorf("origin/HEAD not found in current branch history")
		}

		commit, err = commit.Parent(0)
		if err != nil {
			return nil, err
		}
	}

	return commits, nil
}

func SquashFrom(from, message string) (int, string, error) {
	repo, err := git.PlainOpen(".")
	if err != nil {
		return 0, "", err
	}

	headRef, err := repo.Head()
	if err != nil {
		return 0, "", err
	}

	if !headRef.Name().IsBranch() {
		return 0, "", fmt.Errorf("detached HEAD at %s", headRef.Hash())
	}

	commits, err := getCommitsBetween(from, headRef.Hash().String())
	if err != nil {
		return 0, "", fmt.Errorf("git squash failed: %w", err)
	}
	if len(commits) < 2 {
		return len(commits), "", fmt.Errorf("nothing to squash")
	}
	firstCommit := commits[len(commits)-1]

	headCommit, err := repo.CommitObject(headRef.Hash())
	if err != nil {
		return 0, "", err
	}

	now := time.Now()

	commitMessage := firstCommit.Message
	if message != "" {
		commitMessage = message
	}

	newCommit := &object.Commit{
		Author:       firstCommit.Author,
		Committer:    object.Signature{Name: firstCommit.Author.Name, Email: firstCommit.Author.Email, When: now},
		Message:      commitMessage,
		TreeHash:     headCommit.TreeHash,
		ParentHashes: []plumbing.Hash{plumbing.NewHash(from)},
	}

	obj := repo.Storer.NewEncodedObject()
	if err := newCommit.Encode(obj); err != nil {
		return 0, "", fmt.Errorf("git squash failed: %w", err)
	}

	newHash, err := repo.Storer.SetEncodedObject(obj)
	if err != nil {
		return 0, "", fmt.Errorf("git squash failed: %w", err)
	}

	newRef := plumbing.NewHashReference(headRef.Name(), newHash)

	if err := repo.Storer.SetReference(newRef); err != nil {
		return 0, "", fmt.Errorf("git squash failed: %w", err)
	}

	return len(commits), newHash.String(), nil
}

func Push(auth transport.AuthMethod, force bool) error {
	repo, err := git.PlainOpen(".")
	if err != nil {
		return err
	}

	_, branch, err := GetCurrentBranch()
	if err != nil {
		return err
	}

	refSpec := config.RefSpec(
		"refs/heads/" + branch + ":refs/heads/" + branch,
	)

	if force {
		refSpec = config.RefSpec(
			"+refs/heads/" + branch + ":refs/heads/" + branch,
		)
	}

	err = repo.Push(&git.PushOptions{
		RemoteName: "origin",
		RefSpecs:   []config.RefSpec{refSpec},
		Auth:       auth,
	})

	if errors.Is(err, git.NoErrAlreadyUpToDate) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("git push failed: %w", err)
	}

	return nil
}
