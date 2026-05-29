package git

import (
	"fmt"
	"github.com/Mithweth/git-tools/internal/domain"
	"github.com/go-git/go-git/v5"
	"strings"
)

func GetRepository() (*domain.Repository, error) {
	var provider domain.GitProvider

	repo, err := git.PlainOpen(".")
	if err != nil {
		return nil, err
	}

	remote, err := repo.Remote("origin")
	if err != nil {
		return nil, err
	}

	url := strings.TrimSuffix(remote.Config().URLs[0], ".git")
	host, err := remoteHost(url)
	if err != nil {
		return nil, err
	}
	path, err := remotePath(url)
	if err != nil {
		return nil, err
	}

	switch host {
	case "github.com":
		provider = domain.ProviderGitHub
	case "gitlab.com":
		provider = domain.ProviderGitLab
	default:
		return nil, fmt.Errorf("unsupported git remote url: %s", url)
	}

	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid git remote url: %s", url)
	}
	owner := strings.Join(parts[:len(parts)-1], "/")
	repository := parts[len(parts)-1]

	return &domain.Repository{
		Provider: provider,
		Host:     host,
		URL:      url,
		Owner:    owner,
		Name:     repository,
	}, nil
}

func remoteHost(url string) (string, error) {
	url = strings.TrimPrefix(url, "https://")
	url = strings.TrimPrefix(url, "http://")
	url = strings.TrimPrefix(url, "ssh://")

	if strings.Contains(url, "@") {
		parts := strings.SplitN(url, "@", 2)
		url = parts[1]
	}

	if strings.Contains(url, ":") {
		parts := strings.SplitN(url, ":", 2)
		return parts[0], nil
	}

	if strings.Contains(url, "/") {
		parts := strings.SplitN(url, "/", 2)
		return parts[0], nil
	}

	return "", fmt.Errorf("cannot determine host from %q", url)
}

func remotePath(url string) (string, error) {
	url = strings.TrimPrefix(url, "https://")
	url = strings.TrimPrefix(url, "http://")
	url = strings.TrimPrefix(url, "ssh://")

	if strings.Contains(url, "@") {
		parts := strings.SplitN(url, "@", 2)
		url = parts[1]
	}

	if strings.Contains(url, ":") {
		parts := strings.SplitN(url, ":", 2)
		return parts[1], nil
	}

	if strings.Contains(url, "/") {
		parts := strings.SplitN(url, "/", 2)
		return parts[1], nil
	}

	return "", fmt.Errorf("cannot determine path from %q", url)
}
