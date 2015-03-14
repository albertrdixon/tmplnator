package generator

import (
  "fmt"
  l "github.com/Sirupsen/logrus"
  "github.com/albertrdixon/tmplnator/stack"
  "io/ioutil"
  "os"
  "path/filepath"
  "text/template"
)

type file struct {
  body    *template.Template
  src     string
  name    string
  dir     string
  user    int
  group   int
  mode    os.FileMode
  dirmode os.FileMode
  skip    bool
}

func (f *file) setDir(dir string, args ...interface{}) string {
  for i, a := range args {
    if a == nil {
      args[i] = ""
    }
  }
  f.dir = fmt.Sprintf(dir, args...)
  return ""
}

func (f *file) setName(name string, args ...interface{}) string {
  for i, a := range args {
    if a == nil {
      args[i] = ""
    }
  }
  f.name = fmt.Sprintf(name, args...)
  return ""
}

func (f *file) setUser(uid int) string {
  f.user = uid
  return ""
}

func (f *file) setGroup(gid int) string {
  f.group = gid
  return ""
}

func (f *file) setMode(m os.FileMode) string {
  f.mode = m
  return ""
}

func (f *file) setDirMode(dm os.FileMode) string {
  f.dirmode = dm
  return ""
}

func (f *file) setSkip() string {
  f.skip = true
  return ""
}

func (f *file) Src() string {
  return f.src
}

func (f *file) destination() string {
  return filepath.Join(f.dir, f.name)
}

func parseFiles(dir string, def string) (st *stack.Stack, err error) {
  l.WithField("directory", dir).Info("Parsing files")
  st = stack.NewStack()
  err = filepath.Walk(dir, walkfunc(def, st))
  return
}

func walkfunc(def string, st *stack.Stack) filepath.WalkFunc {
  return func(path string, info os.FileInfo, err error) error {
    ext := filepath.Ext(path)
    if info.Mode().IsRegular() && ext != ".skip" && ext != ".ignore" {
      return parseFile(path, def, st)
    } else {
      l.WithField("path", path).Debug("Skipping")
    }
    return nil
  }
}

func parseFile(path string, def string, st *stack.Stack) (err error) {
  l.WithField("path", path).Debug("Parsing file")

  f := newFile(path, def, filepath.Base(path))
  contents, err := ioutil.ReadFile(path)
  if err != nil {
    return
  }

  t, err := newTemplate(path, f).Parse(string(contents))
  if err != nil {
    return
  }
  f.body = t
  st.Push(f)
  return
}

func newFile(path string, def string, name string) *file {
  return &file{
    src:     path,
    name:    name,
    dir:     def,
    mode:    os.FileMode(0644),
    dirmode: os.FileMode(0755),
    user:    os.Geteuid(),
    group:   os.Getegid(),
    skip:    false,
  }
}

func newTemplate(path string, f *file) *template.Template {
  return template.New(path).Funcs(newFuncMap(f))
}
