package cmd

import (
	"fmt"
	"github.com/bradleyjkemp/git-owners/git"
	"github.com/bradleyjkemp/git-owners/resolver"
	"github.com/bradleyjkemp/git-owners/reviewers"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var PrCmd = &cobra.Command{
	Use:   "pr",
	Short: "Calculates a list of reviewers for a PR",
	Long: `Finds all files that have been modified by commits on this branch
compared to the base-branch.

Resolves a small set of reviewers who can approve this PR
(i.e. passing the output of this into "git owners check" will always return true)

Apart from always being a covering set of Owners, the exact reviewers output
by this command is undefined.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		filesToOwners, err := getFilesToOwnersForPR()
		if err != nil {
			return errors.Wrap(err, "failed to get owners for changed files")
		}
		minimalReviewers := reviewers.SuggestReviewers(filesToOwners)

		if len(minimalReviewers) == 0 {
			fmt.Println("No reviewers needed")
		} else {
			fmt.Println(minimalReviewers)
		}
		return nil
	},
}

func getFilesToOwnersForPR() (map[string][]string, error) {
	baseCommit, err := git.FindBaseCommit(baseBranch)
	if err != nil {
		return nil, err
	}

	changedFiles, err := git.FindChangedFiles(baseCommit)
	if err != nil {
		return nil, err
	}
	filesToOwners := make(map[string][]string)

	for _, file := range changedFiles {
		owners, err := resolver.ResolveOwnersAtCommit(file, false, baseCommit)
		if err != nil {
			return nil, err
		}
		filesToOwners[file] = owners
	}

	return filesToOwners, nil
}

func init() {
	RootCmd.AddCommand(PrCmd)
}
