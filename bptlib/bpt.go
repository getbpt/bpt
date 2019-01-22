package bptlib

import (
	"bufio"
	"github.com/getbpt/bpt/pkg"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"sync"

	"golang.org/x/sync/errgroup"
)

// Bpt the main thing
type Bpt struct {
	r           io.Reader
	parallelism int
	Home        string
}

// New creates a new Bpt instance with the given parameters
func New(home string, r io.Reader, p int) *Bpt {
	return &Bpt{
		r:           r,
		parallelism: p,
		Home:        home,
	}
}

// Get all specified packages and returns the shell content to execute
func (a *Bpt) Get() (result string, err error) {
	var g errgroup.Group
	var lock sync.Mutex
	var shs indexedLines
	var idx int
	sem := make(chan bool, a.parallelism)
	scanner := bufio.NewScanner(a.r)
	for scanner.Scan() {
		l := scanner.Text()
		index := idx
		idx++
		sem <- true
		g.Go(func() error {
			defer func() {
				<-sem
			}()
			l = strings.TrimSpace(l)
			if l == "" || l[0] == '#' {
				return nil
			}
			s, berr := pkg.New(a.Home, l).Get()
			lock.Lock()
			shs = append(shs, indexedLine{idx: index, line: s})
			lock.Unlock()
			return berr
		})
	}
	if err = scanner.Err(); err != nil {
		return
	}
	err = g.Wait()
	return shs.String(), err
}

// Home finds the right home folder to use
func Home() string {
	home := os.Getenv("BPT_HOME")
	if home != "" {
		return home
	}
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return path.Join(path.Join(wd, ".runtime"), "packages")
}
