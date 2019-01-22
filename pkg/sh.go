package pkg

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/getbpt/bpt/project"
)

type shPkg struct {
	Project project.Project
}

func (pkg shPkg) Get() (result string, err error) {
	if err = pkg.Project.Download(); err != nil {
		return result, err
	}
	info, err := os.Stat(pkg.Project.Path())
	if err != nil {
		return "", err
	}
	// it is a file, not a folder, so just return it
	if info.Mode().IsRegular() {
		return "source " + pkg.Project.Path(), nil
	}
	for _, glob := range []string{"*.plugin.sh", "*.sh"} {
		files, err := filepath.Glob(filepath.Join(pkg.Project.Path(), glob))
		if err != nil {
			return result, err
		}
		if files == nil {
			continue
		}
		var lines []string
		for _, file := range files {
			lines = append(lines, "source "+file)
		}
		return strings.Join(lines, "\n"), err
	}

	return result, nil
}
