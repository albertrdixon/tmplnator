package tmplnator

import (
	"fmt"

	l "github.com/Sirupsen/logrus"
	"github.com/spf13/afero"
)

type FileWriter int

func (f FileWriter) String() string {
	return fmt.Sprintf("Writer (%d)", f)
}

func (f FileWriter) writeFiles(done <-chan struct{}, files <-chan *File, fs afero.Fs) <-chan *File {
	out := make(chan *File)
	l.Debugf("START: %v", f)

	go func() {
		defer close(out)
		defer l.Debugf("%v: Exiting", f)
		for file := range files {
			l.Debugf("%v: Working on %v", f, file)
			if file.HasErrs() {
				l.Warnf("%v: %q has parse errors, skipping!", f, file.TemplateName())
				continue
			}
			if err := WriteFile(file, fs); err != nil {
				l.Errorf("%v: Failed to write %q: %v", f, file.Name(), err)
				continue
			}
			select {
			case out <- file:
				l.Debugf("%v: Emitting %v", f, file)
			case <-done:
				return
			}
		}
	}()
	return out
}
