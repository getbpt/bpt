package cmd

import (
	"fmt"
	"github.com/getbpt/bpt/bptlib"
	"github.com/getbpt/bpt/project"
	"github.com/spf13/cobra"
	"runtime"
)

// updateCmd represents the update command
var (
	updateCmd = &cobra.Command{
		Use:   "update [<pkg> ...]",
		Short: "Updates previously downloaded packages",
		RunE: func(cmd *cobra.Command, args []string) error {
			home := bptlib.Home()
			if len(args) < 1 {
				fmt.Printf("Updating all packages in %v...\n", home)
				return project.Update(home, parallelism)
			}
			fmt.Printf("Updating packages in %v...\n", home)
			return project.Update(home, parallelism, args...)
		},
	}
)

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().IntVar(&parallelism, "parallelism", runtime.NumCPU(), "the max amount of tasks to launch in parallel")
}
