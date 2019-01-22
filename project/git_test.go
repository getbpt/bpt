package project

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDownloadAllKinds(t *testing.T) {
	urls := []string{
		"jbaranick/zsh-test_module",
		"http://github.com/caarlos0/ports",
		"http://github.com/caarlos0/ports.git",
		"https://github.com/caarlos0/ports",
		"https://github.com/caarlos0/ports.git",
		"git://github.com/caarlos0/ports.git",
		"https://gitlab.com/caarlos0/test.git",
	}
	for _, url := range urls {
		home := home()
		require.NoError(
			t,
			NewGit(home, url).Download(),
			"Repo "+url+" failed to download",
		)
	}
}

func TestDownloadSubmodules(t *testing.T) {
	var home = home()
	var proj = NewGit(home, "https://github.com/fribmendes/geometry")
	var module = filepath.Join(proj.Path(), "lib/zsh-async")
	require.NoError(t, proj.Download())
	require.NoError(t, proj.Update())
	files, err := ioutil.ReadDir(module)
	require.NoError(t, err)
	require.True(t, len(files) > 1)
}

func TestDownloadAnotherBranch(t *testing.T) {
	home := home()
	require.NoError(t, NewGit(home, "https://github.com/caarlos0/jvm branch:gh-pages").Download())
}

func TestUpdateAnotherBranch(t *testing.T) {
	home := home()
	repo := NewGit(home, "https://github.com/caarlos0/jvm branch:gh-pages")
	require.NoError(t, repo.Download())
	alreadyClonedRepo := NewClonedGit(home, "https-COLON--SLASH--SLASH-github.com-SLASH-caarlos0-SLASH-jvm")
	require.NoError(t, alreadyClonedRepo.Update())
}

func TestUpdateExistentLocalRepo(t *testing.T) {
	home := home()
	repo := NewGit(home, "https://github.com/caarlos0/ports")
	require.NoError(t, repo.Download())
	alreadyClonedRepo := NewClonedGit(home, "https-COLON--SLASH--SLASH-github.com-SLASH-caarlos0-SLASH-ports")
	require.NoError(t, alreadyClonedRepo.Update())
}

func TestUpdateNonExistentLocalRepo(t *testing.T) {
	home := home()
	repo := NewGit(home, "https://github.com/caarlos0/ports")
	require.Error(t, repo.Update())
}

func TestDownloadNonExistentRepo(t *testing.T) {
	home := home()
	repo := NewGit(home, "https://github.com/caarlos0/not-a-real-repo")
	require.Error(t, repo.Download())
}

func TestDownloadMalformedRepo(t *testing.T) {
	home := home()
	repo := NewGit(home, "https://github.com/doesn-not-exist-really branch:also-nope")
	require.Error(t, repo.Download())
}

func TestDownloadMultipleTimes(t *testing.T) {
	home := home()
	repo := NewGit(home, "https://github.com/caarlos0/ports")
	require.NoError(t, repo.Download())
	require.NoError(t, repo.Download())
	require.NoError(t, repo.Update())
}

func TestDownloadFolderNaming(t *testing.T) {
	home := home()
	repo := NewGit(home, "https://github.com/caarlos0/ports")
	require.Equal(
		t,
		home+"/https-COLON--SLASH--SLASH-github.com-SLASH-caarlos0-SLASH-ports",
		repo.Path(),
	)
}

func TestSubFolder(t *testing.T) {
	home := home()
	repo := NewGit(home, "https://github.com/robbyrussell/oh-my-zsh path:plugins/aws")
	require.True(t, strings.HasSuffix(repo.Path(), "plugins/aws"))
}

func TestPath(t *testing.T) {
	home := home()
	repo := NewGit(home, "https://github.com/docker/cli path:contrib/completion/zsh/_docker")
	require.True(t, strings.HasSuffix(repo.Path(), "contrib/completion/zsh/_docker"))
}

func TestMultipleSubFolders(t *testing.T) {
	home := home()
	require.NoError(t, NewGit(home, strings.Join([]string{
		"https://github.com/robbyrussell/oh-my-zsh path:plugins/aws",
		"https://github.com/robbyrussell/oh-my-zsh path:plugins/battery",
	}, "\n")).Download())
}

func home() string {
	home, err := ioutil.TempDir(os.TempDir(), "bpt")
	if err != nil {
		panic(err.Error())
	}
	return home
}
