package tmplnator

import "fmt"

type FileError struct {
	File     string
	Template string
	Msg      string
}

func (f FileError) Error() string {
	return fmt.Sprintf("%s: %s", f.File, f.Msg)
}

func NewFileError(f *File, format string, items ...interface{}) FileError {
	return FileError{f.info.fileName, f.template.Name(), fmt.Sprintf(format, items...)}
}
