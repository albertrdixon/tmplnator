package file

import (
  "bytes"
  "os"
  "path/filepath"
  tmpl "text/template"
)

type mockFile struct {
  name     string
  example  string
  template *tmpl.Template
  fail     map[string]bool
}

func (mf *mockFile) Write(b *bytes.Buffer, data interface{}) (err error) {
  err = mf.template.Execute(b, data)
  return
}

func (mf *mockFile) Read() ([]byte, error) {
  return []byte(mf.example), nil
}

func (mf *mockFile) Template(t *tmpl.Template) {
  mf.template = t
}

func (mf *mockFile) Destination() string {
  return filepath.Join("testing", "example", "path", mf.name)
}

func (mf *mockFile) Src() string {
  return filepath.Join("testing", "example", "path")
}

func (mf *mockFile) DeleteTemplate() error {
  return nil
}

func (mf *mockFile) setDir(d string, args ...interface{}) string {
  return ""
}

func (mf *mockFile) setName(n string, args ...interface{}) string {
  return ""
}

func (mf *mockFile) setUser(uid int) string {
  return ""
}

func (mf *mockFile) setGroup(gid int) string {
  return ""
}

func (mf *mockFile) setMode(fm os.FileMode) string {
  return ""
}

func (mf *mockFile) setDirMode(dm os.FileMode) string {
  return ""
}

func (mf *mockFile) setSkip() string {
  return ""
}

func newMockFile(e string, n string) File {
  return &mockFile{
    example: e,
    name:    n,
  }
}
