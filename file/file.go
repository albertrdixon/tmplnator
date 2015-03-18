package file

import (
  "bytes"
  "os"
  "text/template"
)

// Testing is set to true for running file tests
var Testing bool

// File describes a tmplnator template file
type File interface {
  Write(*bytes.Buffer, interface{}) error
  Read() ([]byte, error)
  Template(*template.Template)
  Destination() string
  Info() Info
  Output() string
  DeleteTemplate() error
  setDir(string, ...interface{}) string
  setName(string, ...interface{}) string
  setUser(int) string
  setGroup(int) string
  setMode(os.FileMode) string
  setDirMode(os.FileMode) string
  setSkip() string
}

// Info objects have all the info for objects that implement File.
type Info struct {
  Src     string
  Name    string
  Dir     string
  User    int
  Group   int
  Mode    os.FileMode
  Dirmode os.FileMode
}

// NewFile returns a File object. If Testing is true underlying struct is
// a mockFile, otherwise it is a templateFile
func NewFile(path string, defaultDir string) File {
  if Testing {
    return newMockFile(path, defaultDir)
  }
  return newTemplateFile(path, defaultDir)
}

func init() {
  Testing = false
}
