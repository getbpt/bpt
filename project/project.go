package project

import (
	"github.com/getbpt/folder"
	"golang.org/x/sync/errgroup"
	"io/ioutil"
)

// Project is basically any kind of project (git, local, svn, bzr, nfs...)
type Project interface {
	Download() error
	Update() error
	Remove() error
	Path() string
	Name() string
	Folder() string
}

// New decides what kind of project it is, based on the given line
func New(home, line string) Project {
	if line[0] == '/' {
		return NewLocal(line)
	}
	return NewGit(home, line)
}

// List all projects in the given folder
func List(home string) (result []string, err error) {
	entries, err := ioutil.ReadDir(home)
	if err != nil {
		return result, err
	}
	for _, entry := range entries {
		if entry.Mode().IsDir() && entry.Name()[0] != '.' {
			result = append(result, entry.Name())
		}
	}
	return result, nil
}

// Update projects in the given folder
func Update(home string, parallelism int, projects ...string) error {
	folders, err := List(home)
	if err != nil {
		return err
	}
	sem := make(chan bool, parallelism)
	var g errgroup.Group
	for _, f := range folders {
		f := f
		sem <- true
		g.Go(func() error {
			defer func() {
				<-sem
			}()
			project := New(home, folder.ToURL(f))
			if projects != nil && len(projects) > 0 {
				for _, p := range projects {
					if p == project.Name() {
						return project.Update()
					}
				}
				return nil
			}
			return project.Update()
		})
	}
	return g.Wait()
}

// Remove projects in the given folder
func Remove(home string, parallelism int, projects ...string) error {
	folders, err := List(home)
	if err != nil {
		return err
	}
	sem := make(chan bool, parallelism)
	var g errgroup.Group
	for _, f := range folders {
		f := f
		sem <- true
		g.Go(func() error {
			defer func() {
				<-sem
			}()
			project := New(home, folder.ToURL(f))
			if projects != nil && len(projects) > 0 {
				for _, p := range projects {
					if p == project.Name() {
						return project.Remove()
					}
				}
				return nil
			}
			return project.Remove()
		})
	}
	return g.Wait()
}
