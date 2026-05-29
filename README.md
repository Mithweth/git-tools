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

## git-create-pr

### Manual installation
```
go mod download
go build -o git-create-pr cmd/git-create-pr/main.go
sudo mv git-create-pr /usr/local/bin/
echo "git config --global alias.create-pr \!/usr/local/bin/git-create-pr" >> ~/.bashrc
```

### Usage

Needs:

* a github.com repository
* a valid GITHUB_TOKEN environment variable

```
git commit -m "my super commit"
git push -u origin mybranch
git create-pr
```

## git-squash

### Manual installation

```
go mod download
go build -o git-squash cmd/git-squash/main.go
sudo mv git-squash /usr/local/bin/
echo "git config --global alias.squash \!/usr/local/bin/git-squash" >> ~/.bashrc
```

### Usage

Needs:

* a git repository
* a valid GITHUB_TOKEN environment variable to push

```
git add foo && git commit -m "my super commit"
git add bar && git commit -m "my other super commit"
git squash -m "my final commit"
```
