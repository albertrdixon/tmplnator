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

func (f FileWriter) writeFiles(files <-chan *File, fs afero.Fs) <-chan *File {
	out := make(chan *File)
	l.Debugf("START: %v", f)

	go func() {
		defer close(out)
		for file := range files {
			l.Debugf("%v: Working on %v.", f, file)
			if file.HasErrs() {
				l.Warnf("%v: %q has parse errors, skipping!", f, file.TemplateName())
				continue
			}
			if err := WriteFile(file, fs); err != nil {
				l.Errorf("%v: Failed to write %q: %v", f, file.Name(), err)
				continue
			}
			out <- file
		}
	}()
	return out
}
