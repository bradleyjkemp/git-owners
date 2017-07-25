package cmd

import (
	"fmt"
	"github.com/bradleyjkemp/git-owners/resolver"
	"github.com/spf13/cobra"
)

var allOwners bool

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all owners of the given files",
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, file := range args {
			owners, err := resolver.ResolveOwners(file, allOwners)
			if err != nil {
				return err
			}
			fmt.Println(file, ":", owners)
		}
		return nil
	},
}

func init() {
	ListCmd.Flags().BoolVarP(&allOwners, "all-owners", "a", false, "Show all owners in the tree, not just the most direct")

	RootCmd.AddCommand(ListCmd)
}
