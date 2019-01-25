package shell

import (
	"bytes"
	"os"
	"text/template"
)

const tmpl = `#!/usr/bin/env bash
bpt () {
  case "$1" in
    (get) eval "$({{ . }} $@ )" || {{ . }} $@ ;;
    (*) {{ . }} $@ ;;
  esac
}
`

// Init returns the shell that should be loaded to bpt to work correctly.
func Init() (string, error) {
	executable, err := os.Executable()
	if err != nil {
		return "", err
	}
	var template = template.Must(template.New("init").Parse(tmpl))
	var out bytes.Buffer
	err = template.Execute(&out, executable)
	return out.String(), err
}
