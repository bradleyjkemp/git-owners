package main

import (
	"flag"
	"fmt"
	"github.com/bradleyjkemp/git-owners/git"
	"github.com/bradleyjkemp/git-owners/resolver"
	"github.com/bradleyjkemp/git-owners/reviewers"
	"os"
)

func main() {
	baseBranch := flag.String("base-branch", "master", "Base branch to compare commits against (default master)")
	allOwners := flag.Bool("a", false, "Resolve all owners up to the root")
	flag.Parse()

	if flag.NArg() == 0 {
		prReviewers(*baseBranch, *allOwners)
	} else {
		for _, file := range flag.Args() {
			owners, err := resolver.ResolveOwners(file, *allOwners)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println(file, ": ", owners)
		}
	}
}

func prReviewers(baseBranch string, allOwners bool) {
	baseCommit, err := git.FindBaseCommit(baseBranch)
	if err != nil {
		fmt.Errorf("%s", err)
	}

	changedFiles, err := git.FindChangedFiles(baseCommit)
	if err != nil {
		fmt.Errorf("%s", err)
	}
	filesToOwners := make(map[string][]string)

	for _, file := range changedFiles {
		owners, err := resolver.ResolveOwnersAtCommit(file, false, baseCommit)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		filesToOwners[file] = owners
	}

	minimalReviewers := reviewers.SuggestReviewers(filesToOwners)

	fmt.Println(minimalReviewers)
}
