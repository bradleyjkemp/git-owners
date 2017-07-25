package cmd

import (
	"fmt"
	"github.com/bradleyjkemp/git-owners/git"
	"github.com/bradleyjkemp/git-owners/resolver"
	"github.com/bradleyjkemp/git-owners/reviewers"
	"github.com/spf13/cobra"
)

// prCmd represents the pr command
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
		baseCommit, err := git.FindBaseCommit(baseBranch)
		if err != nil {
			return err
		}

		changedFiles, err := git.FindChangedFiles(baseCommit)
		if err != nil {
			return err
		}
		filesToOwners := make(map[string][]string)

		for _, file := range changedFiles {
			owners, err := resolver.ResolveOwnersAtCommit(file, false, baseCommit)
			if err != nil {
				return err
			}
			filesToOwners[file] = owners
		}

		minimalReviewers := reviewers.SuggestReviewers(filesToOwners)

		fmt.Println(minimalReviewers)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(PrCmd)
}
