package file

import (
  "bytes"
  "os"
  "text/template"
)

type File interface {
  Write(*bytes.Buffer, interface{}) error
  Read() ([]byte, error)
  Template(*template.Template)
  Destination() string
  Src() string
  DeleteTemplate() error
  setDir(string, ...interface{}) string
  setName(string, ...interface{}) string
  setUser(int) string
  setGroup(int) string
  setMode(os.FileMode) string
  setDirMode(os.FileMode) string
  setSkip() string
}

func newFile(args ...string) File {
  if len(args) == 3 {
    return newTemplateFile(args[0], args[1], args[2])
  }
  return newMockFile(args[0], args[1])
}
