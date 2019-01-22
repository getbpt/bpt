package cmd

import (
	"github.com/getbpt/bpt/version"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the bpt version",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println(version.Print())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
