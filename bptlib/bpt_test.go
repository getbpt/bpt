package bptlib

import (
	"bytes"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBpt(t *testing.T) {
	home := home()
	packages := []string{
		"# comments also are allowed",
		"https://github.com/caarlos0/ports kind:path # comment at the end of the line",
		"https://github.com/caarlos0/jvm kind:path branch:gh-pages",
		"https://github.com/caarlos0/zsh-open-pr kind:sh",
		"",
		"        ",
		"  # trick play",
		"/tmp kind:path",
	}
	sh, err := New(
		home,
		bytes.NewBufferString(strings.Join(packages, "\n")),
		runtime.NumCPU(),
	).Get()
	require.NoError(t, err)
	files, err := ioutil.ReadDir(home)
	require.NoError(t, err)
	require.Len(t, files, 3)
	require.Contains(t, sh, `export PATH="/tmp:$PATH"`)
	require.Contains(t, sh, `export PATH="`+home+`/https-COLON--SLASH--SLASH-github.com-SLASH-caarlos0-SLASH-ports:$PATH"`)
	require.Contains(t, sh, `export PATH="`+home+`/https-COLON--SLASH--SLASH-github.com-SLASH-caarlos0-SLASH-jvm:$PATH"`)
	// nolint: lll
	require.Contains(t, sh, `source `+home+`/https-COLON--SLASH--SLASH-github.com-SLASH-caarlos0-SLASH-zsh-open-pr/git-open-pr.sh`)
}

func TestBptError(t *testing.T) {
	home := home()
	packages := bytes.NewBufferString("invalid-repo")
	sh, err := New(home, packages, runtime.NumCPU()).Get()
	require.Error(t, err)
	require.Empty(t, sh)
}

func TestMultipleRepositories(t *testing.T) {
	home := home()
	packages := []string{
		"# this block is in alphabetic order",
		"https://github.com/caarlos0/git-add-remote kind:path",
		"https://github.com/caarlos0/jvm",
		"https://github.com/caarlos0/ports kind:path",
		"https://github.com/caarlos0/zsh-git-fetch-merge kind:path",
		"https://github.com/caarlos0/zsh-git-sync kind:path",
		"https://github.com/caarlos0/zsh-mkc",
		"https://github.com/caarlos0/zsh-open-pr kind:path",
		"https://github.com/rupa/z",
		"https://github.com/Tarrasch/zsh-bd",
		"https://github.com/wbinglee/zsh-wakatime",
		"https://github.com/zsh-users/zsh-completions",
		"https://github.com/zsh-users/zsh-autosuggestions",
		"",
		"https://github.com/robbyrussell/oh-my-zsh path:plugins/asdf",
		"https://github.com/robbyrussell/oh-my-zsh path:plugins/autoenv",
		"# these should be at last!",
		"https://github.com/sindresorhus/pure",
	}
	sh, err := New(
		home,
		bytes.NewBufferString(strings.Join(packages, "\n")),
		runtime.NumCPU(),
	).Get()
	require.NoError(t, err)
	require.Len(t, strings.Split(sh, "\n"), 15)
}

// BenchmarkDownload-8   	       1	2907868713 ns/op	  480296 B/op	    2996 allocs/op v1
// BenchmarkDownload-8   	       1	2708120385 ns/op	  475904 B/op	    3052 allocs/op v2
func BenchmarkDownload(b *testing.B) {
	var packages = strings.Join([]string{
		"https://github.com/ohmybash/oh-my-bash path:plugins/aws",
		"https://github.com/caarlos0/git-add-remote kind:path",
		"https://github.com/caarlos0/jvm",
		"https://github.com/caarlos0/ports kind:path",
		"",
		"# comment whatever",
		"https://github.com/ohmybash/oh-my-bash path:plugins/battery",
		"https://github.com/rupa/z",
		"",
		"https://github.com/wbinglee/zsh-wakatime",
		"https://github.com/robbyrussell/oh-my-zsh path:plugins/autoenv",
	}, "\n")
	for i := 0; i < b.N; i++ {
		home := home()
		_, err := New(
			home,
			bytes.NewBufferString(packages),
			runtime.NumCPU(),
		).Get()
		require.NoError(b, err)
	}
}

func TestHome(t *testing.T) {
	require.Contains(t, Home(), ".runtime/packages")
}

func TestHomeFromEnvironmentVariable(t *testing.T) {
	require.NoError(t, os.Setenv("BPT_HOME", "/tmp"))
	require.Equal(t, "/tmp", Home())
}

func home() string {
	home, err := ioutil.TempDir(os.TempDir(), "bpt")
	if err != nil {
		panic(err.Error())
	}
	return home
}
