package pkg

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSuccessfullGitPackages(t *testing.T) {
	table := []struct {
		line, result string
	}{
		{
			"https://github.com/caarlos0/jvm",
			"jvm.sh",
		},
		{
			"https://github.com/caarlos0/jvm kind:path",
			"export PATH=\"",
		},
		{
			"https://github.com/caarlos0/jvm kind:path branch:gh-pages",
			"export PATH=\"",
		},
		{
			"https://github.com/caarlos0/jvm kind:dummy",
			"",
		},
		{
			"https://github.com/docker/cli path:contrib/completion/zsh/_docker",
			"contrib/completion/zsh/_docker",
		},
	}
	for _, row := range table {
		row := row
		t.Run(row.line, func(t *testing.T) {
			t.Parallel()
			home := home()
			result, err := New(home, row.line).Get()
			require.Contains(t, result, row.result)
			require.NoError(t, err)
		})
	}
}

func TestShInvalidGitPackage(t *testing.T) {
	home := home()
	_, err := New(home, "does not exist").Get()
	require.Error(t, err)
}

func TestShLocalPackage(t *testing.T) {
	home := home()
	require.NoError(t, ioutil.WriteFile(home+"/a.sh", []byte("echo 9"), 0644))
	result, err := New(home, home).Get()
	require.Contains(t, result, "a.sh")
	require.NoError(t, err)
}

func TestShInvalidLocalPackage(t *testing.T) {
	home := home()
	_, err := New(home, "/asduhasd/asdasda").Get()
	require.Error(t, err)
}

func TestShPackageWithNoShFiles(t *testing.T) {
	home := home()
	_, err := New(home, "https://github.com/getantibody/antibody").Get()
	require.NoError(t, err)
}

func TestPathInvalidLocalPackage(t *testing.T) {
	home := home()
	_, err := New(home, "/asduhasd/asdasda kind:path").Get()
	require.Error(t, err)
}

func TestPathLocalPackage(t *testing.T) {
	home := home()
	require.NoError(t, ioutil.WriteFile(home+"whatever.sh", []byte(""), 0644))
	result, err := New(home, home+" kind:path").Get()
	require.Equal(t, "export PATH=\""+home+":$PATH\"", result)
	require.NoError(t, err)
}

func home() string {
	home, err := ioutil.TempDir(os.TempDir(), "bpt")
	if err != nil {
		panic(err.Error())
	}
	return home
}
