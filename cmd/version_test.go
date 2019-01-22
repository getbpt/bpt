package cmd

import (
	"bytes"
	"strconv"
	"testing"
)

func TestVersionCommand(t *testing.T) {
	output := new(bytes.Buffer)
	rootCmd.SetArgs([]string{"version"})
	rootCmd.SetOutput(output)
	e := rootCmd.Execute()
	if e != nil {
		t.Errorf("Version command failed with: %v", e)
	}
	expected := strconv.Quote(`bpt, version  (branch: , revision: )
  build user:       
  build date:       
  go version:       go1.11.3
`)
	result := strconv.Quote(output.String())
	if result != expected {
		t.Errorf("Version command: want '%v', got '%s'", expected, result)
	}
}
