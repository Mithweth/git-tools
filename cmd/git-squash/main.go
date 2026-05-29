package main

import (
	"fmt"
	"github.com/Mithweth/git-tools/internal/git"
	"github.com/spf13/pflag"
	"os"
)

func Squash(from string, message string) (string, error) {
	var (
		baseHash string
		err      error
	)
	if from == "" {
		baseHash, _, err = git.GetDefaultBranch()
		if err != nil {
			return "", err
		}
	} else {
		baseHash, err = git.ResolveRevision(from)
		if err != nil {
			return "", err
		}
	}
	if commitLength, newHash, errSquash := git.SquashFrom(baseHash, message); errSquash != nil {
		return "", errSquash
	} else {
		return fmt.Sprintf("Squashed %d commits into %s", commitLength, newHash), nil
	}
}

func main() {
	var (
		commitMessage string
		revision      string
		push          bool
		force         bool
	)
	pflag.StringVarP(&commitMessage, "message", "m", "", "Commit message")
	//pflag.BoolVarP(&push, "push", "p", false, "push after squash")
	//pflag.BoolVarP(&force, "force", "f", false, "use force-push instead push")
	pflag.Parse()

	if pflag.NArg() > 0 {
		revision = pflag.Arg(0)
	}
	message, err := Squash(revision, commitMessage)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(message)
	if push {
		if err := git.Push(force); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	}
	os.Exit(0)
}
