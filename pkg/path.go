package pkg

import "github.com/getbpt/bpt/project"

type pathPackage struct {
	Project project.Project
}

func (pkg pathPackage) Get() (result string, err error) {
	if err = pkg.Project.Download(); err != nil {
		return result, err
	}
	return "export PATH=\"" + pkg.Project.Path() + ":$PATH\"", err
}
