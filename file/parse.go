package file

import (
  l "github.com/Sirupsen/logrus"
  "os"
  "path/filepath"
  "text/template"
)

// ParseFiles will recursively parse all the files under dir, returning
// a Queue object with all the files loaded in.
func ParseFiles(dir string, def string) (fq *Queue, err error) {
  l.WithField("directory", dir).Info("Parsing files")
  fq = NewFileQueue()
  err = filepath.Walk(dir, walkfunc(def, fq))
  fq.PopulateQueue()
  return
}

func walkfunc(def string, fq *Queue) filepath.WalkFunc {
  return func(path string, info os.FileInfo, err error) error {
    ext := filepath.Ext(path)
    if info.Mode().IsRegular() && ext != ".skip" && ext != ".ignore" {
      f := NewFile(path, def)
      return ParseFile(f, fq)
    }
    l.WithField("path", path).Debug("Skipping")
    return nil
  }
}

// ParseFile will parse an individual file and put it in the
// Queue
func ParseFile(f File, fq *Queue) (err error) {
  l.WithField("path", f.Info().Src).Debug("Parsing file")

  contents, err := f.Read()
  if err != nil {
    return
  }

  t, err := newTemplate(f).Parse(string(contents))
  if err != nil {
    return
  }
  f.Template(t)
  fq.add(f)
  return
}

func newTemplate(f File) *template.Template {
  return template.New(f.Info().Src).Funcs(newFuncMap(f))
}
