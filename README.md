# git-tools

* git-create-pr: create a pull-request via GitHub API and returns pull request URL to stdout

## Build

```
make
```

## Installation

```
make install-bash
```

Or

```
go mod download
go build -o git-create-pr cmd/git-create-pr/main.go
sudo mv git-create-pr /usr/local/bin/
echo "git config --global alias.create-pr \!/usr/local/bin/git-create-pr" >> ~/.bashrc
```

## Usage

Needs:

* a github.com repository
* a valid GITHUB_TOKEN environment variable

### git-create-pr

```
git commit -m "my super commit"
git push -u origin mybranch
git create-pr
```

Environment variables:

* COMMIT_MESSAGE: Commit message (defaults to the last commit message)
* BRANCH_NAME: Pull request origin branch (default to current branch)
* BASE_BRANCH: Pull request destination branch (default to the repository default branch)
