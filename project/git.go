package project

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/getbpt/folder"
)

type gitProject struct {
	name    string
	Version string
	URL     string
	folder  string
	path    string
	inner   string
}

// NewClonedGit is a git project that was already cloned, so, only Update
// will work here.
func NewClonedGit(home, folderName string) Project {
	folderPath := filepath.Join(home, folderName)
	version, err := branch(folderPath)
	if err != nil {
		version = "master"
	}
	url := folder.ToURL(folderName)
	var name string
	switch {
	case strings.HasPrefix(url, "https://bitbucket.org/"):
		name = strings.TrimPrefix(url, "https://bitbucket.org/")
	case strings.HasPrefix(url, "https://gitlab.com/"):
		fallthrough
	case strings.HasPrefix(url, "https://github.com/"):
		name = strings.TrimPrefix(url, "https://")
	case strings.HasPrefix(url, "http://"):
		fallthrough
	case strings.HasPrefix(url, "https://"):
		fallthrough
	case strings.HasPrefix(url, "git://"):
		fallthrough
	case strings.HasPrefix(url, "ssh://"):
		fallthrough
	case strings.HasPrefix(url, "git@gitlab.com:"):
		fallthrough
	case strings.HasPrefix(url, "git@bitbucket.org:"):
		fallthrough
	case strings.HasPrefix(url, "git@github.com:"):
		name = url
	}
	return gitProject{
		name:    name,
		Version: version,
		URL:     url,
		folder:  folderName,
		path:    folderPath,
	}
}

const (
	branchMarker = "branch:"
	pathMarker   = "path:"
)

// NewGit A git project can be any repository in any given branch. It will
// be downloaded to the provided cwd
func NewGit(cwd, line string) Project {
	version := "master"
	inner := ""
	parts := strings.Split(line, " ")
	for _, part := range parts {
		if strings.HasPrefix(part, branchMarker) {
			version = strings.Replace(part, branchMarker, "", -1)
		}
		if strings.HasPrefix(part, pathMarker) {
			inner = strings.Replace(part, pathMarker, "", -1)
		}
	}
	repo := parts[0]
	url := "https://bitbucket.org/" + repo
	switch {
	case strings.HasPrefix(repo, "http://"):
		fallthrough
	case strings.HasPrefix(repo, "https://"):
		fallthrough
	case strings.HasPrefix(repo, "git://"):
		fallthrough
	case strings.HasPrefix(repo, "ssh://"):
		fallthrough
	case strings.HasPrefix(repo, "git@gitlab.com:"):
		fallthrough
	case strings.HasPrefix(repo, "git@bitbucket.org:"):
		fallthrough
	case strings.HasPrefix(repo, "git@github.com:"):
		url = repo
	case strings.HasPrefix(repo, "gitlab.com/"):
		fallthrough
	case strings.HasPrefix(repo, "bitbucket.org/"):
		fallthrough
	case strings.HasPrefix(repo, "github.com/"):
		url = "https://" + repo
	}
	folder := folder.FromURL(url)
	path := filepath.Join(cwd, folder)
	return gitProject{
		name:    repo,
		Version: version,
		URL:     url,
		folder:  folder,
		path:    path,
		inner:   inner,
	}
}

// nolint: gochecknoglobals
var locks sync.Map

func (g gitProject) Download() error {
	l, _ := locks.LoadOrStore(g.path, &sync.Mutex{})
	lock := l.(*sync.Mutex)
	lock.Lock()
	defer lock.Unlock()
	if _, err := os.Stat(g.path); os.IsNotExist(err) {
		// #nosec
		var cmd = exec.Command("git", "clone",
			"--recursive",
			"--depth", "1",
			"-b", g.Version,
			g.URL,
			g.path)
		cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")

		if bts, err := cmd.CombinedOutput(); err != nil {
			log.Println("git clone failed for", g.URL, string(bts))
			return err
		}
	}
	return nil
}

func (g gitProject) Update() error {
	// #nosec
	if bts, err := exec.Command(
		"git", "-C", g.path, "pull",
		"--recurse-submodules",
		"origin",
		g.Version,
	).CombinedOutput(); err != nil {
		log.Println("git update failed for", g.path, string(bts))
		return err
	}
	return nil
}

func (g gitProject) Remove() error {
	return os.RemoveAll(g.path)
}

func branch(folder string) (string, error) {
	// #nosec
	branch, err := exec.Command(
		"git", "-C", folder, "rev-parse", "--abbrev-ref", "HEAD",
	).Output()
	return strings.Replace(string(branch), "\n", "", -1), err
}

func (g gitProject) Path() string {
	return filepath.Join(g.path, g.inner)
}

func (g gitProject) Name() string {
	return g.name
}

func (g gitProject) Folder() string {
	return g.folder
}
