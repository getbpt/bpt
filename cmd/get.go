package cmd

import (
	"bytes"
	"fmt"
	"github.com/getbpt/bpt/bptlib"
	"github.com/spf13/cobra"
	"io"
	"os"
	"runtime"
	"strings"
)

// getCmd represents the get command
var (
	packageFile string
	getCmd      = &cobra.Command{
		Use:   "get [<pkg>] [<options> ...]",
		Short: "Downloads a pkg and prints its source line",
		RunE: func(cmd *cobra.Command, args []string) error {
			var input io.Reader
			if packageFile != "" {
				var err error
				if input, err = os.Open(packageFile); err != nil {
					return fmt.Errorf("could not open pkg file '%s'", packageFile)
				}
			} else {
				input = bytes.NewBufferString(strings.Join(args, " "))
			}
			sh, err := bptlib.New(bptlib.Home(), input, parallelism).Get()
			if err != nil {
				return err
			}
			fmt.Println(sh)
			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().StringVarP(&packageFile, "requirements", "r", "", "the pkg requirements file")
	getCmd.Flags().IntVar(&parallelism, "parallelism", runtime.NumCPU(), "the max amount of tasks to launch in parallel")
}
