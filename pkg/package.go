package pkg

import (
	"strings"

	"github.com/getbpt/bpt/project"
)

// Package main interface.
type Package interface {
	Get() (result string, err error)
}

// New pkg with at the given home (when apply) and line.
//
// Accepted line formats:
//
// - Local pkg (download and update do nothing):
// 		/home/getbpt/Code/my-local-pkg
// - Bitbucket repo in the owner/repo format:
//		getbpt/bitbucket-repo
// - Git repo in any valid URL form:
//		https://github.com/getbpt/other-github-repo.git
// - Any type of repo, specifying the kind of resource:
//		getbpt/add-to-path-style kind:path
// - Any git repo, specifying a branch:
//		getbpt/versioned-with-branch branch:v1.0 kind:sh
func New(home, line string) Package {
	proj := project.New(home, line)
	switch kind(line) {
	case "binary":
		return binaryPackage{Project: proj}
	case "path":
		return pathPackage{Project: proj}
	case "dummy":
		return dummyPackage{Project: proj}
	default:
		return shPkg{Project: proj}
	}
}

func kind(line string) string {
	for _, part := range strings.Split(line, " ") {
		if strings.HasPrefix(part, "kind:") {
			return strings.Replace(part, "kind:", "", -1)
		}
	}
	return "sh"
}
