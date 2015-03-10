package template

import (
  "crypto/sha1"
  "fmt"
  "github.com/albertrdixon/tmplnator/stack"
  "github.com/oxtoacart/bpool"
  "io/ioutil"
  "os"
  "path/filepath"
  "strings"
  "text/template"
)

type Template struct {
  template *template.Template
  Sha1     string
  Env      map[string]string
  Src      string
}

var (
  files map[string]*File
  env   map[string]string
  bp    *bpool.BufferPool
)

func (t Template) Write() error {
  b := bp.Get()
  defer bp.Put(b)

  err := t.template.Execute(b, t)
  if err != nil {
    return err
  }

  file := files[t.Sha1]
  fmt.Printf("==> Generating %q from %q\n", file.Destination(), t.Src)
  if err := os.MkdirAll(file.DestinationDir(), 0755); err != nil {
    return err
  }
  f, err := os.OpenFile(file.Destination(), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, file.Mode)
  if err != nil {
    return err
  }
  defer f.Close()

  os.Chown(file.Destination(), file.User, file.Group)
  _, err = f.Write(b.Bytes())
  if err != nil {
    return err
  }
  return nil
}

func ParseDirectory(dir string, prefix string) (*stack.Stack, error) {
  st := stack.NewStack()
  return st, filepath.Walk(dir, walkfunc(st, prefix))
}

func walkfunc(tmplStack *stack.Stack, prefix string) filepath.WalkFunc {
  return func(path string, info os.FileInfo, err error) error {
    if info.Mode().IsRegular() {
      return ParseTemplate(path, prefix, tmplStack)
    }
    return nil
  }
}

func ParseTemplate(path string, prefix string, st *stack.Stack) error {
  base := filepath.Base(path)
  t := new(Template)
  t.Src = path
  t.Env = env
  contents, err := ioutil.ReadFile(path)
  if err != nil {
    return err
  }

  t.Sha1 = fmt.Sprintf("%x", sha1.New().Sum(contents))
  template, err := newTemplate(path).Parse(string(contents))
  if err != nil {
    return err
  }

  t.template = template
  st.Push(t)
  files[t.Sha1] = newFile(prefix, base)
  return nil
}

func newTemplate(path string) *template.Template {
  return template.New(path).Funcs(newFuncMap())
}

func envMap() map[string]string {
  env := make(map[string]string, len(os.Environ()))
  for _, val := range os.Environ() {
    index := strings.Index(val, "=")
    env[val[:index]] = val[index+1:]
  }
  return env
}

func iface(list []string) []interface{} {
  vals := make([]interface{}, len(list))
  for i, v := range list {
    vals[i] = v
  }
  return vals
}

func init() {
  bp = bpool.NewBufferPool(48)
  env = envMap()
  files = make(map[string]*File)
}
