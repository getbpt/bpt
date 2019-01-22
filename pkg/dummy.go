package pkg

import "github.com/getbpt/bpt/project"

type dummyPackage struct {
	Project project.Project
}

func (pkg dummyPackage) Get() (result string, err error) {
	err = pkg.Project.Download()
	return result, err
}
