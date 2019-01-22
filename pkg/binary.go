package pkg

import (
	"github.com/getbpt/bpt/project"
	"os"
	"path/filepath"
	"runtime"
)

type binaryPackage struct {
	Project project.Project
}

func (pkg binaryPackage) Get() (result string, err error) {
	if err = pkg.Project.Download(); err != nil {
		return result, err
	}
	var platformPath = filepath.Join(pkg.Project.Path(), runtime.GOOS)
	if _, err = os.Stat(platformPath); err != nil {
		return "", err
	}
	return "export PATH=\"" + platformPath + ":$PATH\"", err
}
