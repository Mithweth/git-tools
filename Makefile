GO ?= go
REMOVE ?= rm
INSTALLBIN ?= install

ifeq ($(PREFIX),)
    PREFIX := /usr/local
endif

default: build

build: fmt
	@$(GO) mod download
	@$(GO) build -o git-create-pr cmd/git-create-pr/main.go
	@$(GO) build -o git-squash cmd/git-squash/main.go

fmt:
	@$(GO) fmt ./...

vet:
	@$(GO) vet ./...

test: fmt vet
	@$(GO) test ./... -coverprofile=cover.out

coverage:
	@$(GO) tool cover -func=cover.out

clean:
	@$(REMOVE) -f git-create-pr cover.out git-squash

install: 
	sudo $(INSTALLBIN) -d $(PREFIX)/bin/
	sudo $(INSTALLBIN) -m 755 git-create-pr $(PREFIX)/bin/
	sudo $(INSTALLBIN) -m 755 git-squash $(PREFIX)/bin/

uninstall:
	sudo $(REMOVE) -f $(PREFIX)/bin/git-create-pr
	sudo $(REMOVE) -f $(PREFIX)/bin/git-squash
	sed -i.bak "/alias.create-pr/d" ~/.bashrc
	sed -i.bak "/alias.git-squash/d" ~/.bashrc

install-bash: install
	echo "git config --global alias.create-pr \!$(PREFIX)/bin/git-create-pr" >> ~/.bashrc
	echo "git config --global alias.squash \!$(PREFIX)/bin/git-squash" >> ~/.bashrc
