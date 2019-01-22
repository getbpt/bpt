package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	parallelism int
	osExit      = os.Exit

	// rootCmd represents the base command when called without any sub-commands
	rootCmd = &cobra.Command{
		Use:   "bpt",
		Short: "Bash Package Tool",
		Long: `bpt provides a simple way to declaratively retrieve shell scripts, binaries, etc. 
and expose them to the current shell.

To use, there are two steps to perform in a script:
  1. Initialize bpt: eval "$(bpt init)"
  2. Get pkg: bpt get org/repo`,
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		osExit(1)
	}
}
