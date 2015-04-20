package tmplnator

import (
	"fmt"

	l "github.com/Sirupsen/logrus"
	"github.com/spf13/afero"
)

type FileParser int

func (f FileParser) String() string {
	return fmt.Sprintf("Parser (%d)", f)
}

func (f FileParser) parseTemplates(fs afero.Fs, files ...*File) <-chan *File {
	out := make(chan *File)

	l.Debugf("START: %v", f)
	go func() {
		defer close(out)
		for _, file := range files {
			l.Debugf("%v: Parsing %q.", f, file.TemplateName())
			if err := ParseTemplate(file, srcFs); err != nil {
				l.Errorf("%v: Failed to parse %q: %v", f, file.TemplateName(), err)
				continue
			}
			out <- file
		}
	}()
	return out
}
