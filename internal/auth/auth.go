package auth

import (
	"fmt"
	"github.com/Mithweth/git-tools/internal/domain"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/kevinburke/ssh_config"
	"os"
	"path/filepath"
	"strings"
)

func expandHome(path string) string {
	if strings.HasPrefix(path, "~/") {
		return filepath.Join(os.Getenv("HOME"), path[2:])
	}
	return path
}

func sshAuth(host string) (transport.AuthMethod, error) {
	if envVar := os.Getenv("GIT_SSH_KEY"); envVar != "" {
		auth, err := ssh.NewPublicKeysFromFile("git", envVar, "")
		if err != nil {
			return nil, err
		}
		return auth, nil
	}

	identity := ssh_config.Get(host, "IdentityFile")
	if identity == "" {
		return ssh.NewSSHAgentAuth("git")
	}

	auth, err := ssh.NewPublicKeysFromFile("git", expandHome(identity), "")
	if err != nil {
		return nil, err
	}

	return auth, nil
}

func httpsAuth(host string) (transport.AuthMethod, error) {
	var token string

	switch host {
	case "github.com":
		token = os.Getenv("GITHUB_TOKEN")
	case "gitlab.com":
		token = os.Getenv("GITLAB_TOKEN")
	}

	if token == "" {
		return nil, fmt.Errorf("missing token for HTTPS remote %s", host)
	}

	return &http.BasicAuth{
		Username: "git",
		Password: token,
	}, nil
}

func GetAuth(r *domain.Repository) (transport.AuthMethod, error) {
	if strings.HasPrefix(r.URL, "git@") || strings.HasPrefix(r.URL, "ssh://") {
		return sshAuth(r.Host)
	}

	if strings.HasPrefix(r.URL, "https://") {
		return httpsAuth(r.Host)
	}

	return nil, fmt.Errorf("unsupported remote url: %s", r.URL)
}
