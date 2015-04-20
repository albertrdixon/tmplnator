package tmplnator

import (
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
	srcFs   afero.Fs
	destFs  afero.Fs

	// Backend is the Key/Value store for Data
	Backend backend.Backend

	// ForceTemp will force all templates to be written to TmpDir if true
	ForceTemp = false

	// TmpDir is the tmpdir we will use
	TmpDir = filepath.Join(os.TempDir(), "T2")

	cpus = runtime.NumCPU()
	pkg  = "tmplnator"
)

func InitFs(srcInMem, destInMem bool) {
	mem := &afero.MemMapFs{}
	os := &afero.OsFs{}

	if srcFs == nil {
		srcFs = os
		if srcInMem {
			srcFs = mem
		}
	}

	if destFs == nil {
		destFs = os
		if destInMem {
			destFs = mem
		}
	}
}

func ClearFs() {
	srcFs, destFs = nil, nil
}

// Generate is the main entrypoint. It will parse any templates found under dir.
func Generate(dir string) []*File {
	l.Infof("Starting generator. Template directory: %q", dir)
	var rawFiles []*File
	var finishedFiles []*File
	data = NewData(Backend)
	defer func() { data = nil }()

	if err := destFs.MkdirAll(TmpDir, os.FileMode(0777)); err != nil {
		l.Warnf("Could not create tmpdir %q: %v", TmpDir, err)
		TmpDir = ""
	}

	rawFiles = dirRead(dir, srcFs)
	if len(rawFiles) < 1 {
		l.Errorf("Found no templates under %q!", dir)
		return rawFiles
	}

	l.Infof("Found %d files to parse.", len(rawFiles))
	p := FileParser(0)
	pfc := p.parseTemplates(srcFs, rawFiles...)
	var wfc []<-chan *File
	for i := 0; i < cpus; i++ {
		w := FileWriter(i)
		wfc = append(wfc, w.writeFiles(pfc, destFs))
	}

	for file := range merge(wfc...) {
		finishedFiles = append(finishedFiles, file)
	}

	if len(finishedFiles) < len(rawFiles) {
		n := len(rawFiles) - len(finishedFiles)
		l.Errorf("%d files failed to be parsed or written!", n)
	}
	l.Infof("Generated %d files.", len(finishedFiles))

	return finishedFiles
}

func dirRead(root string, fs afero.Fs) []*File {
	var files []*File
	l.Debugf("Reading dir %s", root)

	items, err := afero.ReadDir(root, fs)
	if err != nil {
		l.Errorf("Unable to read dir %q: %v", root, err)
		return files
	}

	for _, item := range items {
		path := filepath.Join(root, item.Name())
		l.Debugf("Reading item %q", path)
		if item.IsDir() {
			files = append(files, dirRead(path, fs)...)
			continue
		}
		if !strings.Contains(item.Name(), "ignore") && !strings.Contains(item.Name(), "skip") {
			l.Debugf("Adding %q to file queue.", item.Name())
			files = append(files, NewFile(path))
		}
	}
	return files
}

func merge(cs ...<-chan *File) <-chan *File {
	var wg sync.WaitGroup
	out := make(chan *File)

	output := func(n int, c <-chan *File) {
		for f := range c {
			l.Debugf("gatherer (%d): Got Output: %v", n, f)
			out <- f
		}
		l.Debug("Done")
		wg.Done()
	}
	l.Debugf("WaitGroup: %d", len(cs))
	wg.Add(len(cs))
	for i, c := range cs {
		l.Debugf("START: Result Gatherer (%d)", i)
		go output(i, c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
