package file

import (
  "bytes"
  "fmt"
  "os"
  "path/filepath"
  tmpl "text/template"
)

type mockFile struct {
  dir      string
  dirmode  os.FileMode
  example  string
  group    int
  mode     os.FileMode
  name     string
  output   string
  template *tmpl.Template
  user     int
}

func (mf *mockFile) Write(b *bytes.Buffer, data interface{}) (err error) {
  err = mf.template.Execute(b, data)
  mf.output = b.String()
  return
}

func (mf *mockFile) Read() ([]byte, error) {
  return []byte(mf.example), nil
}

func (mf *mockFile) Template(t *tmpl.Template) {
  mf.template = t
}

func (mf *mockFile) Destination() string {
  return filepath.Join(mf.dir, mf.name)
}

func (mf *mockFile) Info() Info {
  return Info{
    Name:    mf.name,
    Dir:     mf.dir,
    User:    mf.user,
    Group:   mf.group,
    Mode:    mf.mode,
    Dirmode: mf.dirmode,
  }
}

func (mf *mockFile) Output() string {
  return mf.output
}

func (mf *mockFile) DeleteTemplate() error {
  return nil
}

func (mf *mockFile) setDir(d string, args ...interface{}) string {
  for i, a := range args {
    if a == nil {
      args[i] = ""
    }
  }
  mf.dir = fmt.Sprintf(d, args...)
  return ""
}

func (mf *mockFile) setName(n string, args ...interface{}) string {
  for i, a := range args {
    if a == nil {
      args[i] = ""
    }
  }
  mf.name = fmt.Sprintf(n, args...)
  return ""
}

func (mf *mockFile) setUser(uid int) string {
  mf.user = uid
  return ""
}

func (mf *mockFile) setGroup(gid int) string {
  mf.group = gid
  return ""
}

func (mf *mockFile) setMode(fm os.FileMode) string {
  mf.mode = fm
  return ""
}

func (mf *mockFile) setDirMode(dm os.FileMode) string {
  mf.dirmode = dm
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
