package tmplnator

import (
	"os"
	"path/filepath"
	"sync"

	l "github.com/Sirupsen/logrus"
	"github.com/spf13/afero"
)

func initFs(srcInMem, destInMem bool) {
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

// Generate is the main entrypoint. It will parse any templates found under dir.
func generate(dir string) []*File {
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

	done := make(chan struct{})
	defer close(done)

	pfc := FileParser(0).parseTemplates(done, srcFs, rawFiles...)
	var wfc []<-chan *File
	for i := 0; i < cpus; i++ {
		wfc = append(wfc, FileWriter(i).writeFiles(done, pfc, destFs))
	}

	for file := range merge(done, wfc...) {
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
		if isTemplate(item.Name()) {
			l.Debugf("Adding %q to file queue.", item.Name())
			files = append(files, NewFile(path))
		} else {
			l.Debugf("Skipping %q.", item.Name())
		}
	}
	return files
}

func merge(done <-chan struct{}, cs ...<-chan *File) <-chan *File {
	var wg sync.WaitGroup
	out := make(chan *File)

	l.Debugf("WaitGroup size: %d", len(cs))
	wg.Add(len(cs))
	for i, c := range cs {
		go Gatherer(i).gatherOutput(&wg, done, out, c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
