package cmd

import (
	"fmt"
	"github.com/getbpt/bpt/bptlib"
	"github.com/spf13/cobra"
)

// homeCmd represents the home command
var homeCmd = &cobra.Command{
	Use:   "home",
	Short: "Prints the root directory where packages are kept",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(bptlib.Home())
	},
}

func init() {
	rootCmd.AddCommand(homeCmd)
}
