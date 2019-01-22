package cmd

import (
	"fmt"
	"github.com/getbpt/bpt/bptlib"
	"github.com/getbpt/bpt/project"
	"github.com/spf13/cobra"
	"runtime"
)

// purgeCmd represents the purge command
var (
	purgeCmd = &cobra.Command{
		Use:   "purge [<pkg> ...]",
		Short: "Purges previously downloaded packages",
		RunE: func(cmd *cobra.Command, args []string) error {
			home := bptlib.Home()
			if len(args) < 1 {
				fmt.Printf("Removing all packages in %v...\n", home)
				return project.Remove(home, parallelism)
			}
			fmt.Printf("Removing packages in %v...\n", home)
			return project.Remove(home, parallelism, args...)
		},
	}
)

func init() {
	rootCmd.AddCommand(purgeCmd)
	purgeCmd.Flags().IntVar(&parallelism, "parallelism", runtime.NumCPU(), "the max amount of tasks to launch in parallel")
}
