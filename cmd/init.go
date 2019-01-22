package cmd

import (
	"fmt"
	"github.com/getbpt/bpt/shell"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Configure the shell environment for bpt",
	RunE: func(cmd *cobra.Command, args []string) error {
		sh, err := shell.Init()
		fmt.Println(sh)
		return err
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
