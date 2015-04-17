package tmplnator

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	l "github.com/Sirupsen/logrus"
	"github.com/albertrdixon/tmplnator/backend"
	"github.com/spf13/afero"
)

var (
	errChan chan error
	data    *Data

	// Backend is the Key/Value store for Data
	Backend backend.Backend

	// ForceTemp will force all templates to be written to TmpDir if true
	ForceTemp = false

	// TmpDir is the tmpdir we will use
	TmpDir = filepath.Join(os.TempDir(), "T2")

	cpus = runtime.NumCPU()
	pkg  = "tmplnator"
)

// Generator holds references to the filesystem. It will parse templates
// found in srcFs and write generated files to destFs.
type Generator struct {
	srcFs  afero.Fs
	destFs afero.Fs
}

func NewGenerator(srcInMem, destInMem bool) *Generator {
	mem := &afero.MemMapFs{}
	os := &afero.OsFs{}

	var src afero.Fs
	var dest afero.Fs
	src, dest = os, os
	if srcInMem {
		src = mem
	}
	if destInMem {
		dest = mem
	}
	return &Generator{src, dest}
}

// Generate is the main entrypoint. It will parse any templates found under dir.
func (g *Generator) Generate(dir string) ([]string, []error) {
	if err := g.destFs.MkdirAll(TmpDir, os.FileMode(0777)); err != nil {
		TmpDir = ""
	}

	errChan = make(chan error, 10)
	data = NewData(Backend)

	l.WithField("directory", dir).Info("Starting generator.")
	paths := dirRead(dir, g.srcFs)
	l.Infof("Found %d files to parse.", len(paths))
	templates := parseTemplates(g.srcFs, paths...)
	var fcs []<-chan string
	for i := 0; i < cpus; i++ {
		fcs = append(fcs, writeFiles(i, templates, g.destFs))
	}

	var files []string
	var errs []error
	for file := range merge(fcs...) {
		files = append(files, file)
	}
	close(errChan)
	if len(errChan) > 0 {
		for err := range errChan {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		l.Warnf("Saw %d errors!", len(errs))
	}
	l.Infof("Wrote %d files.", len(files))
	return files, errs
}

func parseTemplates(fs afero.Fs, paths ...string) <-chan *Template {
	out := make(chan *Template)

	l.Debug("START: Parser")
	go func() {
		defer close(out)
		for _, path := range paths {
			l.Debugf("Parsing %s", path)
			file, err := fs.Open(path)
			if err != nil {
				l.Errorf("Could not read template %s: %v", path, err)
				continue
			}
			t := NewTemplate(path)
			if err := t.Parse(file); err == nil {
				out <- t
			} else {
				errChan <- newGeneratorError("parser", "%v", err)
			}
			file.Close()
		}
	}()
	return out
}

func writeFiles(n int, templates <-chan *Template, fs afero.Fs) <-chan string {
	thread := fmt.Sprintf("writer (%d)", n)
	defer func() {
		if r := recover(); r != nil {
			l.WithFields(l.Fields{
				"writer":  n,
				"package": pkg,
			}).Errorf("Recover from panic: %v", r)
			errChan <- newGeneratorError(thread, "%v", r)
		}
	}()

	out := make(chan string)
	l.Debugf("START: Writer (%d)", n)
	go func() {
		defer close(out)
		for t := range templates {
			if t.template == nil {
				continue
			}
			l.Debugf("(%d) Working on %s", n, t.Info.Name)
			buf := new(bytes.Buffer)
			if err := t.template.Execute(buf, data); err != nil {
				errChan <- newGeneratorError(thread, "Could not exec template %s: %v", t.Info.Source, err)
				continue
			}

			if t.Skip() {
				continue
			}

			if _, err := fs.Stat(t.Info.Dir); err != nil {
				if err := fs.MkdirAll(t.Info.Dir, t.Info.Dirmode); err != nil {
					errChan <- newGeneratorError(thread, "Could not mkdir %s: %v", t.Info.Dir, err)
					continue
				}
			}

			file, err := fs.Create(t.Output())
			if err != nil {
				errChan <- newGeneratorError(thread, "Could not open %s: %v", t.Output(), err)
				continue
			}

			if _, err = buf.WriteTo(file); err != nil {
				errChan <- newGeneratorError(thread, "Could not write to %s: %v", t.Output(), err)
				continue
			}
			fs.Chmod(t.Output(), t.Info.Mode)
			out <- file.Name()
		}
	}()
	return out
}

func dirRead(root string, fs afero.Fs) []string {
	var files []string
	l.Debugf("Reading dir %s", root)

	items, err := afero.ReadDir(root, fs)
	if err != nil {
		errChan <- newGeneratorError("main", "%v", err)
		return files
	}

	for _, item := range items {
		l.Debugf("Reading item %s", item)
		if item.IsDir() {
			files = append(files, dirRead(filepath.Join(root, item.Name()), fs)...)
			continue
		}
		if !strings.Contains(item.Name(), "ignore") && !strings.Contains(item.Name(), "skip") {
			files = append(files, filepath.Join(root, item.Name()))
		}
	}
	return files
}

func merge(cs ...<-chan string) <-chan string {
	var wg sync.WaitGroup
	out := make(chan string)

	output := func(c <-chan string) {
		for f := range c {
			l.Debugf("Got Output: %v", f)
			out <- f
		}
		l.Debug("Done")
		wg.Done()
	}
	l.Debugf("WaitGroup: %d", len(cs))
	wg.Add(len(cs))
	for i, c := range cs {
		l.Debugf("Get Output (%d)", i)
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
