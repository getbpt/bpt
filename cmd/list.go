package cmd

import (
	"fmt"
	"github.com/getbpt/bpt/bptlib"
	"github.com/getbpt/bpt/project"
	"github.com/getbpt/folder"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"text/tabwriter"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Prints the currently installed packages",
	RunE: func(cmd *cobra.Command, args []string) error {
		home := bptlib.Home()
		projects, err := project.List(home)
		if err != nil {
			return err
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 1, 4, ' ', tabwriter.TabIndent)
		for _, b := range projects {
			_, _ = fmt.Fprintf(w, "%s\t%s\n", folder.ToURL(b), filepath.Join(home, b))
		}
		return w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
