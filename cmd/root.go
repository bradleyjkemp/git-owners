package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var baseBranch string

var RootCmd = &cobra.Command{
	Use:   "git-owners",
	Args:  cobra.ArbitraryArgs,
	Short: "A tool for finding owners and reviewers for files",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			return ListCmd.RunE(ListCmd, args)
		} else {
			return PrCmd.RunE(PrCmd, nil)
		}

	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&baseBranch, "base-branch", "b", "master", "The branch this PR is being merged into")
}
