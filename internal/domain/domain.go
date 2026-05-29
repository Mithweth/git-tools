package domain

type GitProvider string

const (
	ProviderGitHub GitProvider = "github"
	ProviderGitLab GitProvider = "gitlab"
)

type Repository struct {
	Provider GitProvider
	Host     string
	Owner    string
	Name     string
	URL      string
}
