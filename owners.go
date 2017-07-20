package main

import (
	"flag"
	"fmt"
	"github.com/bradleyjkemp/git-owners/git"
	"github.com/bradleyjkemp/git-owners/resolver"
	"github.com/bradleyjkemp/git-owners/reviewers"
	"io"
	"os"
)

var stdout io.Writer = os.Stdout
var stderr io.Writer = os.Stderr

type cliFlags struct {
	baseBranch string
	allOwners  bool
	args       []string
}

func main() {
	baseBranch := flag.String("base-branch", "master", "Base branch to compare commits against (default master)")
	allOwners := flag.Bool("a", false, "Resolve all owners up to the root")
	flag.Parse()

	owners(cliFlags{
		*baseBranch,
		*allOwners,
		flag.Args(),
	})
}

func owners(flags cliFlags) {
	if len(flags.args) == 0 {
		prReviewers(flags)
	} else {
		fileOwners(flags)
	}
}

func fileOwners(flags cliFlags) {
	for _, file := range flags.args {
		owners, err := resolver.ResolveOwners(file, flags.allOwners)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Fprintln(stdout, file, ":", owners)
	}
}

func prReviewers(flags cliFlags) {
	baseCommit, err := git.FindBaseCommit(flags.baseBranch)
	if err != nil {
		fmt.Fprintf(stderr, "%s\n", err)
		os.Exit(1)
	}

	changedFiles, err := git.FindChangedFiles(baseCommit)
	if err != nil {
		fmt.Fprintf(stderr, "%s\n", err)
		os.Exit(1)
	}
	filesToOwners := make(map[string][]string)

	for _, file := range changedFiles {
		owners, err := resolver.ResolveOwnersAtCommit(file, flags.allOwners, baseCommit)
		if err != nil {
			fmt.Fprintln(stderr, err)
			os.Exit(1)
		}
		filesToOwners[file] = owners
	}

	minimalReviewers := reviewers.SuggestReviewers(filesToOwners)

	fmt.Fprintln(stdout, minimalReviewers)
}
