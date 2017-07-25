package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// checkCmd represents the check command
var CheckCmd = &cobra.Command{
	Use:   "check reviewer1 reviewer2",
	Short: "Check a PR has been accepted by all necessary owners",
	Long:  `Given a list of reviewers this checks that every file modified on this branch has at least one owner who has accepted the PR`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("unimplemented")
	},
}

func init() {
	RootCmd.AddCommand(CheckCmd)
}
